package auth

import "errors"

var ErrDuplicateEmail = errors.New("duplicate email")
var ErrDuplicateUsername = errors.New("duplicate username")
var ErrDuplicatePhone = errors.New("duplicate phone")
var ErrInvalidEmail = errors.New("invalid email")
var ErrInvalidUsername = errors.New("invalid username")
var ErrInvalidPassword = errors.New("invalid password")
var ErrRequiredEmail = errors.New("email field required")
var ErrRequiredUsername = errors.New("username field required")
var ErrMissingField = errors.New("missing field")
var ErrDatabaseJob = errors.New("database job failed")
