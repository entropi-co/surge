package api

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"surge/internal/auth"
	"surge/internal/utilities"
)

func (a *SurgeAPI) EndpointSignUpWithCredentials(w http.ResponseWriter, r *http.Request) error {
	body, err := utilities.GetBodyJson[SignUpWithCredentialsRequest](r)
	if err != nil {
		return err
	}

	var validationErrors validator.ValidationErrors

	queries := a.queries

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*body.Password), bcrypt.DefaultCost)
	if err != nil {
		return InternalServerError("failed to hash password")
	}
	hashedPasswordString := string(hashedPassword[:])

	createdUser, err := auth.CreateUser(queries, r.Context(), a.config, auth.CreateUserOptions{
		Phone:    body.Phone,
		Email:    body.Email,
		Username: body.Username,
		Password: &hashedPasswordString,
		Metadata: auth.UserMetadata{
			Avatar:    body.Metadata.Avatar,
			FirstName: body.Metadata.FirstName,
			LastName:  body.Metadata.LastName,
			Birthdate: body.Metadata.Birthdate,
		},
	})
	if err != nil {
		switch {
		case errors.Is(auth.ErrMissingField, err):
		case errors.Is(auth.ErrRequiredUsername, err):
		case errors.Is(auth.ErrRequiredEmail, err):
			return BadRequestError(ErrorCodeMissingField, err.Error())
		case errors.Is(auth.ErrInvalidEmail, err):
		case errors.Is(auth.ErrInvalidUsername, err):
		case errors.Is(auth.ErrInvalidPassword, err):
			return BadRequestError(ErrorCodeInvalidField, err.Error())
		case errors.Is(auth.ErrDuplicateEmail, err):
		case errors.Is(auth.ErrDuplicateUsername, err):
		case errors.Is(auth.ErrDuplicatePhone, err):
			return ConflictError(err.Error())
		case errors.Is(auth.ErrDatabaseJob, err):
			return InternalServerError("failed to do database action")
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

		return InternalServerError("unknown error during creating user")
	}

	return writeResponseJSON(w, http.StatusOK, NewUserResponse(createdUser))
}
