package api

type ErrorCode = string

const (
	// ErrorCodeUnknown should not be used directly, it only indicates a failure in the error handling system in such a way that an error code was not assigned properly.
	ErrorCodeUnknown ErrorCode = "unknown"

	// ErrorCodeUnexpectedFailure signals an unexpected failure such as a 500 Internal Server Error.
	ErrorCodeUnexpectedFailure ErrorCode = "unexpected_failure"
	ErrorCodeDatabaseFailure   ErrorCode = "database_failure"

	ErrorCodeBadJSON      ErrorCode = "bad_json"
	ErrorCodeInvalidJSON  ErrorCode = "invalid_json"
	ErrorCodeMissingField ErrorCode = "missing_field"
	ErrorCodeInvalidField ErrorCode = "invalid_field"

	ErrorCodeInvalidCredentials ErrorCode = "invalid_credentials"

	ErrorCodeConflict ErrorCode = "conflict"

	ErrorCodeInvalidGrantType  ErrorCode = "invalid_grant_type"
	ErrorCodeDisabledGrantType ErrorCode = "disabled_grant_type"
)
