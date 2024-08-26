package api

type ErrorCode = string

const (
	// ErrorCodeUnknown should not be used directly, it only indicates a failure in the error handling system in such a way that an error code was not assigned properly.
	ErrorCodeUnknown ErrorCode = "unknown"

	// ErrorCodeUnexpectedFailure signals an unexpected failure such as a 500 Internal Server Error.
	ErrorCodeUnexpectedFailure ErrorCode = "unexpected_failure"
	ErrorCodeDatabaseFailure   ErrorCode = "database_failure"

	ErrorCodeInvalidJSON  ErrorCode = "invalid_json"
	ErrorCodeMissingField ErrorCode = "missing_field"
	ErrorCodeInvalidField ErrorCode = "invalid_field"

	ErrorCodeInvalidCredentials ErrorCode = "invalid_credentials"

	ErrorCodeConflict ErrorCode = "conflict"

	ErrorCodeInvalidGrantType  ErrorCode = "invalid_grant_type"
	ErrorCodeDisabledGrantType ErrorCode = "disabled_grant_type"

	ErrorCodeInvalidProviderType ErrorCode = "invalid_provider_type"

	ErrorCodeBadOAuth2State    ErrorCode = "bad_oauth2_state"
	ErrorCodeBadOAuth2Callback ErrorCode = "bad_oauth2_callback"

	ErrorCodeProviderOAuth2Unsupported ErrorCode = "provider_oauth2_unsupported"

	ErrorCodeUserNotFound ErrorCode = "user_not_found"

	ErrorCodeRefreshNotFoundToken ErrorCode = "refresh_token_not_found"
	ErrorCodeRefreshTokenRevoked  ErrorCode = "refresh_token_revoked"

	ErrorCodeNoAuthorization ErrorCode = "no_authorization"
	ErrorCodeBadJWT          ErrorCode = "bad_jwt"
)
