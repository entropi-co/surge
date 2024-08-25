package api

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

// HTTPError is an error with a message and an HTTP status code.
type HTTPError struct {
	HTTPStatus int    `json:"code"`
	ErrorCode  string `json:"error_code,omitempty"`
	Message    string `json:"message,omitempty"`
	Details    any    `json:"details,omitempty"`
	ErrorID    string `json:"error_id,omitempty"`
}

// Builder Start

// HTTPErrorBuilder is builder for HTTPError
type HTTPErrorBuilder struct {
	httpError *HTTPError
}

func NewBuilder() *HTTPErrorBuilder {
	return &HTTPErrorBuilder{
		httpError: &HTTPError{},
	}
}

func (b *HTTPErrorBuilder) SetStatus(status int) *HTTPErrorBuilder {
	b.httpError.HTTPStatus = status
	return b
}

func (b *HTTPErrorBuilder) SetErrorCode(code ErrorCode) *HTTPErrorBuilder {
	b.httpError.ErrorCode = code
	return b
}

func (b *HTTPErrorBuilder) SetMessage(fmtString string, args ...interface{}) *HTTPErrorBuilder {
	b.httpError.Message = fmt.Sprintf(fmtString, args...)
	return b
}

func (b *HTTPErrorBuilder) SetDetails(details any) *HTTPErrorBuilder {
	b.httpError.Details = details
	return b
}

func (b *HTTPErrorBuilder) SetErrorID(id string) *HTTPErrorBuilder {
	b.httpError.ErrorID = id
	return b
}

func (b *HTTPErrorBuilder) UseRequest(r *http.Request) *HTTPErrorBuilder {
	return b.SetErrorID(middleware.GetReqID(r.Context()))
}

func (b *HTTPErrorBuilder) Build() *HTTPError {
	return b.httpError
}

// Builder End

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%d: %s", e.HTTPStatus, e.Message)
}

func (e *HTTPError) Is(target error) bool {
	return e.Error() == target.Error()
}

// Cause returns the root cause error
func (e *HTTPError) Cause() error {
	return e
}

func UnauthorizedError(errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return NewHTTPError(http.StatusUnauthorized, errorCode, fmtString, args...)
}

func BadRequestError(errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return NewHTTPError(http.StatusBadRequest, errorCode, fmtString, args...)
}

func InternalServerError(fmtString string, args ...interface{}) *HTTPError {
	return NewHTTPError(http.StatusInternalServerError, ErrorCodeUnexpectedFailure, fmtString, args...)
}

func NotFoundError(errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return NewHTTPError(http.StatusNotFound, errorCode, fmtString, args...)
}

func ForbiddenError(errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return NewHTTPError(http.StatusForbidden, errorCode, fmtString, args...)
}

func UnprocessableEntityError(errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return NewHTTPError(http.StatusUnprocessableEntity, errorCode, fmtString, args...)
}

func TooManyRequestsError(errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return NewHTTPError(http.StatusTooManyRequests, errorCode, fmtString, args...)
}

func ConflictError(fmtString string, args ...interface{}) *HTTPError {
	return NewHTTPError(http.StatusConflict, ErrorCodeConflict, fmtString, args...)
}

func NewHTTPError(httpStatus int, errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return &HTTPError{
		HTTPStatus: httpStatus,
		ErrorCode:  errorCode,
		Message:    fmt.Sprintf(fmtString, args...),
	}
}
