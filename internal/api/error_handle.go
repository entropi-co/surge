package api

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
)

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
