package api

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"surge/internal/schema"
	"time"
)

// SurgeAPI is exposed API for Surge
type SurgeAPI struct {
	httpHandler http.Handler
	version     *string
	tx          schema.DBTX
}

// NewSurgeAPI Creates a new SurgeAPI instance
func NewSurgeAPI() SurgeAPI {
	api := SurgeAPI{
		version: nil,
	}

	api.httpHandler = api.createHttpHandler()

	return api
}

// ListenAndServe starts the REST API with httpHandler
func (a *SurgeAPI) ListenAndServe(ctx context.Context, hostAndPort string) {
	baseCtx, cancel := context.WithCancel(context.Background())

	log := logrus.WithField("component", "api")

	server := &http.Server{
		Addr:              hostAndPort,
		Handler:           a.httpHandler,
		ReadHeaderTimeout: 2 * time.Second, // to mitigate a Slowloris attack
		BaseContext: func(net.Listener) context.Context {
			return baseCtx
		},
	}

	cleanupWaitGroup.Add(1)
	go func() {
		defer cleanupWaitGroup.Done()

		<-ctx.Done()

		defer cancel() // close baseContext

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Minute)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, context.Canceled) {
			log.WithError(err).Error("shutdown failed")
		}
	}()

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.WithError(err).Fatal("http server listen failed")
	}
}
