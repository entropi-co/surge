package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"net/http"
	"surge/internal/utilities"
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

	router.Get("/health", a.EndpointHealth)

	router.Route("/v1", func(router *SurgeAPIRouter) {
		router.Route("/sign_up", func(router *SurgeAPIRouter) {
			router.Post("/credentials", a.EndpointSignUpWithCredentials)
		})

		router.Post("/token", a.EndpointToken)
	})

	var totalRouteNodes = 0
	var totalRouteEndpoints = 0
	utilities.Walk(router.chi.Routes(), func(route chi.Route) []chi.Route {
		if route.SubRoutes == nil {
			totalRouteNodes++
			totalRouteEndpoints++
			return []chi.Route{}
		}

		totalRouteNodes++
		return route.SubRoutes.Routes()
	})

	logger.Infof("Created router with %d nodes, %d endpoints", totalRouteNodes, totalRouteEndpoints)

	return router
}
