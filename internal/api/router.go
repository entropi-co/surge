package api

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"surge/internal/utilities"
)

type SurgeAPIRouter struct {
	chi chi.Router
}

func NewSurgeAPIRouter() *SurgeAPIRouter {
	return &SurgeAPIRouter{chi.NewRouter()}
}

func (r *SurgeAPIRouter) Route(pattern string, fn func(*SurgeAPIRouter)) {
	r.chi.Route(pattern, func(c chi.Router) {
		fn(&SurgeAPIRouter{c})
	})
}

func (r *SurgeAPIRouter) Get(pattern string, fn surgeAPIHandler) {
	r.chi.Get(pattern, handler(fn))
}
func (r *SurgeAPIRouter) Post(pattern string, fn surgeAPIHandler) {
	r.chi.Post(pattern, handler(fn))
}
func (r *SurgeAPIRouter) Put(pattern string, fn surgeAPIHandler) {
	r.chi.Put(pattern, handler(fn))
}
func (r *SurgeAPIRouter) Delete(pattern string, fn surgeAPIHandler) {
	r.chi.Delete(pattern, handler(fn))
}

func (r *SurgeAPIRouter) With(fn middlewareHandler) *SurgeAPIRouter {
	c := r.chi.With(wrapAsMiddleware(fn))
	return &SurgeAPIRouter{c}
}

func (r *SurgeAPIRouter) WithBypass(fn func(next http.Handler) http.Handler) *SurgeAPIRouter {
	c := r.chi.With(fn)
	return &SurgeAPIRouter{c}
}

func (r *SurgeAPIRouter) Use(fn middlewareHandler) {
	r.chi.Use(wrapAsMiddleware(fn))
}

func (r *SurgeAPIRouter) UseBypass(fn func(next http.Handler) http.Handler) {
	r.chi.Use(fn)
}
func (r *SurgeAPIRouter) UseRequestLogging() {
	r.Use(func(w http.ResponseWriter, r *http.Request) (context.Context, error) {
		logrus.
			WithContext(r.Context()).
			WithField("agent", r.UserAgent()).
			Infof("Request: %s %s", r.Method, r.URL.Path)

		return r.Context(), nil
	})
}

func (r *SurgeAPIRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.chi.ServeHTTP(w, req)
}

func (r *SurgeAPIRouter) CountNodes() (int, int) {
	var totalRouteNodes = 0
	var totalRouteEndpoints = 0
	utilities.Walk(r.chi.Routes(), func(route chi.Route) []chi.Route {
		if route.SubRoutes == nil {
			totalRouteNodes++
			totalRouteEndpoints++
			return []chi.Route{}
		}

		totalRouteNodes++
		return route.SubRoutes.Routes()
	})

	return totalRouteNodes, totalRouteEndpoints
}

type surgeAPIHandler func(w http.ResponseWriter, r *http.Request) error

func handler(fn surgeAPIHandler) http.HandlerFunc {
	return fn.serve
}

func (h surgeAPIHandler) serve(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		HandleResponseError(err, w, r)
	}
}

type middlewareHandler func(w http.ResponseWriter, r *http.Request) (context.Context, error)

func (m middlewareHandler) handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.serve(next, w, r)
	})
}

func (m middlewareHandler) serve(next http.Handler, w http.ResponseWriter, r *http.Request) {
	ctx, err := m(w, r)
	if err != nil {
		HandleResponseError(err, w, r)
		return
	}
	if ctx != nil {
		r = r.WithContext(ctx)
	}
	next.ServeHTTP(w, r)
}

func wrapAsMiddleware(fn middlewareHandler) func(http.Handler) http.Handler {
	return fn.handler
}
