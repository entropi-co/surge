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
	Message    string `json:"message"`
	ErrorID    string `json:"error_id,omitempty"`
}

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
