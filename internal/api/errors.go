package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
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

func unauthorizedError(errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusUnauthorized, errorCode, fmtString, args...)
}

func badRequestError(errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusBadRequest, errorCode, fmtString, args...)
}

func internalServerError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusInternalServerError, ErrorCodeUnexpectedFailure, fmtString, args...)
}

func notFoundError(errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusNotFound, errorCode, fmtString, args...)
}

func forbiddenError(errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusForbidden, errorCode, fmtString, args...)
}

func unprocessableEntityError(errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusUnprocessableEntity, errorCode, fmtString, args...)
}

func tooManyRequestsError(errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusTooManyRequests, errorCode, fmtString, args...)
}

func conflictError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusConflict, ErrorCodeConflict, fmtString, args...)
}

func httpError(httpStatus int, errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return &HTTPError{
		HTTPStatus: httpStatus,
		ErrorCode:  errorCode,
		Message:    fmt.Sprintf(fmtString, args...),
	}
}
func HandleResponseError(err error, w http.ResponseWriter, r *http.Request) {
	log := logrus.New()
	requestId := middleware.GetReqID(r.Context())

	var e *HTTPError
	switch {
	case errors.As(err, &e):
		switch {
		case e.HTTPStatus >= http.StatusInternalServerError:
			e.ErrorID = requestId
			// this will get us the stack trace too
			log.WithError(e.Cause()).Error(e.Error())
		case e.HTTPStatus == http.StatusTooManyRequests:
			log.WithError(e.Cause()).Warn(e.Error())
		default:
			log.WithError(e.Cause()).Info(e.Error())
		}

		if e.ErrorCode == "" {
			if e.HTTPStatus == http.StatusInternalServerError {
				e.ErrorCode = ErrorCodeUnexpectedFailure
			} else {
				e.ErrorCode = ErrorCodeUnknown
			}
		}

		if jsonErr := writeResponseJSON(w, e.HTTPStatus, e); jsonErr != nil && jsonErr != context.DeadlineExceeded {
			log.WithError(jsonErr).Warn("Failed to send JSON on ResponseWriter")
		}

		return
	}

	log.WithError(err).Errorf("Unhandled server error: %s", err.Error())

	httpError := HTTPError{
		HTTPStatus: http.StatusInternalServerError,
		ErrorCode:  ErrorCodeUnexpectedFailure,
		Message:    "Unexpected failure, please check server logs for more information",
		ErrorID:    requestId,
	}

	if jsonErr := writeResponseJSON(w, http.StatusInternalServerError, httpError); jsonErr != nil && !errors.Is(jsonErr, context.DeadlineExceeded) {
		log.WithError(jsonErr).Warn("Failed to send JSON on ResponseWriter")
	}
}
