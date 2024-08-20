package api

import (
	"net/http"
	"surge/internal/utilities"
)

type signInWithCredentialsBody struct {
	username *string
	email    *string
	password *string
}

func (a *SurgeAPI) EndpointSignInWithCredentials(w http.ResponseWriter, r *http.Request) error {
	body, err := utilities.GetBodyJson[signInWithCredentialsBody](r)
	if err != nil {
		return err
	}

	if body.username == nil && body.email == nil {
		return badRequestError(ErrorCodeBadJSON, "neither username nor email was provided")
	}
	if body.username != nil && body.email != nil {
		return badRequestError(ErrorCodeBadJSON, "both username and email was provided")
	}
	if body.password == nil {
		return badRequestError(ErrorCodeBadJSON, "password was not provided")
	}

	return httpError(http.StatusUnauthorized, ErrorCodeUserNotFound, "user with the credentials was not found")
}
