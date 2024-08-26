package api

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
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
	logger := logrus.WithField("component", "router")

	router := NewSurgeAPIRouter()
	router.UseBypass(middleware.RequestID)

	if a.config.Logging.EnableRequest {
		router.UseRequestLogging()
		logger.Infoln("Enabled request logging")
	}

	router.Get("/health", a.EndpointHealth)
	router.Get("/.well-known/jwks.json", a.EndpointJwks)

	router.Route("/v1", func(router *SurgeAPIRouter) {
		router.Route("/sign_up", func(router *SurgeAPIRouter) {
			router.Post("/credentials", a.EndpointSignUpWithCredentials)
		})

		router.Post("/token", a.EndpointToken)

		router.Route("/external", func(router *SurgeAPIRouter) {
			router.Get("/", a.EndpointExternal)

			router.Route("/callback", func(router *SurgeAPIRouter) {
				router.Use(a.loadOAuth2StateToContextMiddleware)

				router.Get("/", a.EndpointExternalCallback)
				router.Post("/", a.EndpointExternalCallback)
			})
		})

		router.Route("/user", func(router *SurgeAPIRouter) {
			router.Use(a.useAuthentication)

			router.Get("/", a.EndpointUser)
			// TODO: Add update user route (POST|PUT /user)
		})
	})

	totalRouteNodes, totalRouteEndpoints := router.CountNodes()
	logger.
		WithField("nodes", totalRouteNodes).
		WithField("endpoints", totalRouteEndpoints).
		Infoln("Created router")

	return router
}
