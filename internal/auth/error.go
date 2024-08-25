package auth

import "errors"

var (
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrDuplicatePhone    = errors.New("duplicate phone")

	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidPassword = errors.New("invalid password")

	ErrRequiredEmail    = errors.New("email field required")
	ErrRequiredUsername = errors.New("username field required")
	ErrRequiredPhone    = errors.New("username field required")

	ErrMissingField = errors.New("missing field")
	ErrDatabaseJob  = errors.New("database job failed")
)
