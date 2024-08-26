package api

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"surge/internal/auth"
	"surge/internal/conf"
	"surge/internal/schema"
	"surge/internal/storage"
	"surge/internal/utilities"
	"time"
)

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	Email    *string `json:"email"`
	Username *string `json:"username"`
}

type TokenGrantType = string

const (
	TokenGrantTypeCredentials TokenGrantType = "credentials"
	TokenGrantTypeRefresh     TokenGrantType = "refresh"
)

type tokenCredentialsGrantTypeRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type tokenRefreshGrantTypeRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// EndpointToken endpoint used to log in a user and respond with accessToken
func (a *SurgeAPI) EndpointToken(w http.ResponseWriter, r *http.Request) error {
	grantType := r.URL.Query().Get("grant_type")

	switch grantType {
	case TokenGrantTypeCredentials:
		return a.tokenCredentialsGrantFlow(w, r)
	case TokenGrantTypeRefresh:
		return a.tokenRefreshGrantFlow(w, r)
	default:
		return BadRequestError(ErrorCodeInvalidGrantType, "invalid grant type '%s'", grantType)
	}
}

// tokenCredentialsGrantFlow processes grant flow with username/email and password
func (a *SurgeAPI) tokenCredentialsGrantFlow(w http.ResponseWriter, r *http.Request) error {
	authorizationErr := UnauthorizedError(ErrorCodeInvalidCredentials, "failed to find matching user with the credentials")

	body, err := utilities.GetBodyJson[tokenCredentialsGrantTypeRequest](r)
	if err != nil {
		return err
	}

	// Email and Username field is provided at the same time
	if utilities.CountNotNil([]*string{body.Email, body.Username}) != 1 {
		return BadRequestError(ErrorCodeInvalidJSON, "email or username can't be provided at the same time")
	}

	var user *schema.AuthUser

	if body.Email != nil {
		// Abort if email auth is disabled
		if a.config.Auth.DisableEmailAuth {
			return UnprocessableEntityError(ErrorCodeDisabledGrantType, "email authentication is disabled")
		}

		user, err = a.queries.GetUserByEmail(r.Context(), *body.Email)
	} else if body.Username != nil {
		// Abort if username auth is disabled
		if a.config.Auth.DisableUsernameAuth {
			return UnprocessableEntityError(ErrorCodeDisabledGrantType, "username authentication is disabled")
		}

		user, err = a.queries.GetUserByUsername(r.Context(), *body.Username)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authorizationErr
		}
		return NewHTTPError(http.StatusInternalServerError, ErrorCodeDatabaseFailure, "unexpected database failure")
	}

	if !auth.AuthenticateUser(user, *body.Password) {
		return authorizationErr
	}

	token, err := a.issueToken(r.Context(), user)
	if err != nil {
		return err
	}

	return writeResponseJSON(w, http.StatusOK, token)
}

func (a *SurgeAPI) tokenRefreshGrantFlow(w http.ResponseWriter, r *http.Request) error {
	body, err := utilities.GetBodyJson[tokenRefreshGrantTypeRequest](r)
	if err != nil {
		return err
	}

	if body.RefreshToken == "" {
		return BadRequestError(ErrorCodeInvalidField, "refresh_token is empty or missing")
	}

	refreshToken, err := a.queries.GetRefreshToken(r.Context(), body.RefreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NotFoundError(ErrorCodeRefreshNotFoundToken, "failed to find refresh token")
		}
		return err
	}

	if refreshToken.Revoked {
		return ForbiddenError(ErrorCodeRefreshTokenRevoked, "refresh token was revoked")
	}

	user, err := a.queries.GetUser(r.Context(), refreshToken.UserID.UUID)
	if err != nil {
		return err
	}

	var response *AccessTokenResponse

	// Revoke and issue new refresh token
	err = a.Transaction(r.Context(), func(tx *sql.Tx, queries *schema.Queries) error {
		if err := queries.RevokeRefreshToken(r.Context(), refreshToken.ID); err != nil {
			return err
		}

		response, err = a.issueToken(r.Context(), user)
		return err
	})
	if err != nil {
		return err
	}

	return writeResponseJSON(w, http.StatusOK, response)
}

func (a *SurgeAPI) issueToken(ctx context.Context, user *schema.AuthUser) (*AccessTokenResponse, error) {
	logger := logrus.WithContext(ctx).WithField("user", user.ID)

	accessTokenString, expiresAt, err := a.generateAccessToken(user)
	if err != nil {
		var httpErr *HTTPError
		if errors.As(err, &httpErr) {
			return nil, httpErr
		}

		logger.WithError(err).Errorln("failed to generate access accessToken")
		return nil, InternalServerError("failed to generate access accessToken")
	}
	refreshToken, err := a.generateRefreshToken(ctx, a.queries, user)
	if err != nil {
		return nil, InternalServerError("failed to create refresh accessToken")
	}

	return &AccessTokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken.Token.String,
		ExpiresIn:    a.config.JWT.ExpiresAfter,
		ExpiresAt:    expiresAt,
		User:         NewUserResponse(user),
	}, nil
}

func (a *SurgeAPI) generateRefreshToken(ctx context.Context, q *schema.Queries, user *schema.AuthUser) (*schema.AuthRefreshToken, error) {
	logger := logrus.WithContext(ctx).WithField("user", user.ID)

	token, err := q.CreateRefreshToken(ctx, schema.CreateRefreshTokenParams{
		UserID:  uuid.NullUUID{UUID: user.ID, Valid: user != nil},
		Token:   storage.NewString(utilities.SecureToken()),
		Revoked: false,
	})
	if err != nil {
		logger.Errorf("Failed to create refresh accessToken")
		return nil, err
	}

	return token, nil
}

// generateAccessToken generates accessToken with configured JWKs in configuration and returns (accessToken, expiresAt, error)
func (a *SurgeAPI) generateAccessToken(user *schema.AuthUser) (string, int64, error) {
	//logger := logrus.WithField("user", user.ID).WithField("where", "access_token_generation")

	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(time.Second * time.Duration(a.config.JWT.ExpiresAfter))

	claims := AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		Email:    storage.NullStringToPointer(user.Email),
		Username: storage.NullStringToPointer(user.Username),
	}

	// Acquire signing JWK
	signingKey, err := a.config.JWT.GetSigningJwk()
	if err != nil {
		return "", 0, err
	}

	signingMethod := conf.GetJwkCompatibleAlgorithm(signingKey)

	// Create accessToken with claims
	token := jwt.NewWithClaims(signingMethod, claims)

	jwt.MarshalSingleStringAsArray = false
	// Acquire raw signing key from JWK
	rawSigningKey, err := conf.GetSigningKeyFromJwk(signingKey)
	if err != nil {
		return "", 0, err
	}

	// Sign accessToken
	signedToken, err := token.SignedString(rawSigningKey)
	if err != nil {
		return "", 0, err
	}

	return signedToken, expiresAt.Unix(), nil
}
