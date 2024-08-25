package api

import (
	"github.com/lestrrat-go/jwx/v2/jwk"
	"net/http"
)

// EndpointJwks exposed at /.well-known/jwks.json, is used for providing public jwks to first party services for verifying the accessToken
func (a *SurgeAPI) EndpointJwks(w http.ResponseWriter, r *http.Request) error {
	res := JwksResponse{
		Keys: []jwk.Key{},
	}

	for i := range a.config.JWT.Keys {
		keyPair := a.config.JWT.Keys[i]

		// Skip if public key is not present
		if keyPair.PublicKey == nil {
			continue
		}

		res.Keys = append(res.Keys, keyPair.PublicKey)
	}

	w.Header().Set("Cache-Control", "public, max-age=600")
	return writeResponsePrettyJSON(w, http.StatusOK, res)
}
