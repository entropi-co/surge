package api

import (
	"net/http"
	"surge/internal/utilities"
)

type HealthCheckResponse struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

var defaultVersion = "development build"

// HealthCheck endpoint indicates if the gotrue api service is available
func (a *SurgeAPI) EndpointHealth(w http.ResponseWriter, r *http.Request) error {
	return writeResponseJSON(w, http.StatusOK, HealthCheckResponse{
		Version: *utilities.Coalesce(a.version, &defaultVersion),
		Name:    "Surge API",
	})
}
