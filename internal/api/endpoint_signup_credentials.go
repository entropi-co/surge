package api

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"surge/internal/auth"
	"surge/internal/utilities"
	"time"
)

type signUpWithCredentialsBody struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type signUpResponseBody struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *SurgeAPI) EndpointSignUpWithCredentials(w http.ResponseWriter, r *http.Request) error {
	body, err := utilities.GetBodyJson[signUpWithCredentialsBody](r)
	if err != nil {
		return err
	}

	var validationErrors validator.ValidationErrors

	requestId := middleware.GetReqID(r.Context())
	log := logrus.WithContext(r.Context()).WithField("req", requestId)
	queries := a.queries

	log.Infof("Creating user...")

	createdUser, err := auth.CreateUser(queries, r.Context(), a.config, auth.CreateUserOptions{
		Email:    body.Email,
		Username: body.Username,
		Password: body.Password,
	})
	if err != nil {
		switch {
		case errors.Is(auth.ErrMissingField, err):
		case errors.Is(auth.ErrRequiredUsername, err):
		case errors.Is(auth.ErrRequiredEmail, err):
			return badRequestError(ErrorCodeMissingField, err.Error())
		case errors.Is(auth.ErrInvalidEmail, err):
		case errors.Is(auth.ErrInvalidUsername, err):
		case errors.Is(auth.ErrInvalidPassword, err):
			return badRequestError(ErrorCodeInvalidField, err.Error())
		case errors.Is(auth.ErrDuplicateEmail, err):
		case errors.Is(auth.ErrDuplicateUsername, err):
		case errors.Is(auth.ErrDuplicatePhone, err):
			return conflictError(err.Error())
		case errors.Is(auth.ErrDatabaseJob, err):
			return internalServerError("failed to do database action")
		case errors.As(err, &validationErrors):
			httpErr := NewBuilder().
				UseRequest(r).
				SetStatus(http.StatusBadRequest).
				SetErrorCode(ErrorCodeInvalidJSON).
				SetDetails(utilities.Map(
					validationErrors,
					func(t validator.FieldError) any {
						return map[string]any{
							"tag":       t.Tag(),
							"namespace": t.Namespace(),
							"field":     t.Field(),
							"error":     fmt.Sprintf("failed validation for field '%s' on the '%s' tag", t.Field(), t.Tag()),
						}
					},
				)).
				Build()
			return httpErr
		}

		return internalServerError("unknown error during creating user")
	}

	return writeResponseJSON(w, http.StatusOK,
		signUpResponseBody{
			Id:        createdUser.ID,
			CreatedAt: createdUser.CreatedAt,
		},
	)
}
