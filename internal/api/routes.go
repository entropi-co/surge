package api

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"net/http"
)

func (a *SurgeAPI) createHttpHandler() http.Handler {
	corsHandler := cors.New(cors.Options{
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
	})

	return corsHandler.Handler(a.createRouter())
}

func (a *SurgeAPI) createRouter() *SurgeAPIRouter {
	router := NewSurgeAPIRouter()
	router.UseBypass(middleware.RequestID)

	router.Get("/health", a.EndpointHealth)

	router.Route("/v1", func(router *SurgeAPIRouter) {
		router.Route("/sign_up", func(router *SurgeAPIRouter) {
			router.Post("/credentials", a.EndpointSignUpWithCredentials)
		})
	})

	return router
}
