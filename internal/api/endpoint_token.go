package api

import (
	"database/sql"
	"errors"
	"net/http"
	"surge/internal/auth"
	"surge/internal/schema"
	"surge/internal/storage"
	"surge/internal/utilities"
)

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

	return nil
}

func (a *SurgeAPI) issueToken(user *schema.AuthUser) {
	// TODO: Token issue logic
}
