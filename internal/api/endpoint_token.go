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
)

type tokenCredentialsGrantTypeBody struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

// EndpointToken endpoint used to log in a user and respond with token
func (a *SurgeAPI) EndpointToken(w http.ResponseWriter, r *http.Request) error {
	grantType := r.URL.Query().Get("grant_type")

	switch grantType {
	case TokenGrantTypeCredentials:
		return a.tokenCredentialsGrantFlow(w, r)
	default:
		return badRequestError(ErrorCodeInvalidGrantType, "invalid grant type '%s'", grantType)
	}
}

// tokenCredentialsGrantFlow processes grant flow with username/email and password
func (a *SurgeAPI) tokenCredentialsGrantFlow(w http.ResponseWriter, r *http.Request) error {
	authorizationErr := unauthorizedError(ErrorCodeInvalidCredentials, "failed to find matching user with the credentials")

	body, err := utilities.GetBodyJson[tokenCredentialsGrantTypeBody](r)
	if err != nil {
		return err
	}

	// Email and Username field is provided at the same time
	if utilities.CountNotNil([]*string{body.Email, body.Username}) != 1 {
		return badRequestError(ErrorCodeInvalidJSON, "email or username can't be provided at the same time")
	}

	var user schema.AuthUser

	if body.Email != nil {
		// Abort if email auth is disabled
		if a.config.Auth.DisableEmailAuth {
			return unprocessableEntityError(ErrorCodeDisabledGrantType, "email authentication is disabled")
		}

		user, err = a.queries.GetUserByEmail(r.Context(), storage.NewString(*body.Email))
	} else if body.Username != nil {
		// Abort if username auth is disabled
		if a.config.Auth.DisableUsernameAuth {
			return unprocessableEntityError(ErrorCodeDisabledGrantType, "username authentication is disabled")
		}

		user, err = a.queries.GetUserByUsername(r.Context(), storage.NewString(*body.Username))
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authorizationErr
		}
		return httpError(http.StatusInternalServerError, ErrorCodeDatabaseFailure, "unexpected database failure")
	}

	if !auth.AuthenticateUser(&user, *body.Password) {
		return authorizationErr
	}

	token, err := a.issueToken(r.Context(), &user)
	if err != nil {
		return err
	}

	return writeResponseJSON(w, http.StatusOK, token)
}

func (a *SurgeAPI) issueToken(ctx context.Context, user *schema.AuthUser) (*AccessTokenResponse, error) {
	logger := logrus.WithContext(ctx).WithField("user", user.ID)

	accessTokenString, expiresAt, err := a.generateAccessToken(user)
	if err != nil {
		if httpErr, ok := err.(*HTTPError); ok {
			return nil, httpErr
		}

		logger.WithError(err).Errorln("failed to generate access token")
		return nil, internalServerError("failed to generate access token")
	}
	refreshToken, err := a.generateRefreshToken(ctx, user)
	if err != nil {
		return nil, internalServerError("failed to create refresh token")
	}

	return &AccessTokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken.Token.String,
		ExpiresIn:    a.config.JWT.ExpiresAfter,
		ExpiresAt:    expiresAt,
		User:         NewUserResponse(user),
	}, nil
}

func (a SurgeAPI) generateRefreshToken(ctx context.Context, user *schema.AuthUser) (*schema.AuthRefreshToken, error) {
	logger := logrus.WithContext(ctx).WithField("user", user.ID)

	token, err := a.queries.CreateRefreshToken(ctx, schema.CreateRefreshTokenParams{
		UserID:  uuid.NullUUID{UUID: user.ID, Valid: user != nil},
		Token:   storage.NewString(utilities.SecureToken()),
		Revoked: false,
	})
	if err != nil {
		logger.Errorf("Failed to create refresh token")
		return nil, err
	}

	return &token, nil
}

// generateAccessToken generates token with configured JWKs in configuration and returns (token, expiresAt, error)
func (a SurgeAPI) generateAccessToken(user *schema.AuthUser) (string, int64, error) {
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

	signingKey, err := a.config.JWT.GetSigningJwk()
	if err != nil {
		return "", 0, err
	}

	signingMethod := conf.GetJwkCompatibleAlgorithm(signingKey)

	token := jwt.NewWithClaims(signingMethod, claims)

	rawSigningKey, err := conf.GetSigningKeyFromJwk(signingKey)
	if err != nil {
		return "", 0, err
	}

	signedToken, err := token.SignedString(rawSigningKey)
	if err != nil {
		return "", 0, err
	}

	return signedToken, expiresAt.Unix(), nil
}
