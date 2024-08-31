package api

import (
	"net/http"
	"net/url"
	"surge/internal/conf"
)

func GetRequestReferrer(r *http.Request, config *conf.SurgeConfigurations) string {
	// try get redirect url from query or post data first
	reqref := getRedirectTo(r)
	if IsRedirectURLValid(config, reqref) {
		return reqref
	}

	// instead try referrer header value
	reqref = r.Referer()
	if IsRedirectURLValid(config, reqref) {
		return reqref
	}

	return config.ServiceURL
}

func IsRedirectURLValid(config *conf.SurgeConfigurations, redirectURL string) bool {
	if redirectURL == "" {
		return false
	}

	base, berr := url.Parse(config.ServiceURL)
	refurl, rerr := url.Parse(redirectURL)

	// As long as the referrer came from the site, we will redirect back there
	if berr == nil && rerr == nil && base.Hostname() == refurl.Hostname() {
		return true
	}

	// For case when user came from mobile app or other permitted resource - redirect back
	for _, pattern := range config.URIAllowListMap {
		if pattern.Match(redirectURL) {
			return true
		}
	}

	return false
}

func getRedirectTo(r *http.Request) string {
	reqref := r.Header.Get("redirect_to")
	if reqref != "" {
		return reqref
	}

	return r.URL.Query().Get("redirect_to")
}
