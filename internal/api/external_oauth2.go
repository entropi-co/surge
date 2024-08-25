package api

import (
	"context"
	"fmt"
	"net/http"
	"surge/internal/api/provider"
)

// OAuth2ProviderData contains the userData and accessToken returned by the oauth provider
type OAuth2ProviderData struct {
	userData     *provider.UserData
	accessToken  string
	refreshToken string
	code         string
}

func (a *SurgeAPI) handleOAuthCallback(r *http.Request) (*OAuth2ProviderData, error) {
	ctx := r.Context()
	providerType := ctx.Value(contextExternalProviderTypeKey).(string)

	response, err := a.oauth2Callback(r, providerType)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *SurgeAPI) loadOAuth2StateToContextMiddleware(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	state := r.URL.Query().Get("state")
	if state == "" {
		return nil, BadRequestError(ErrorCodeBadOAuth2Callback, "OAuth state parameter missing")
	}

	ctx := r.Context()
	return a.loadExternalStateToContext(ctx, state)
}

func (a *SurgeAPI) oauth2Callback(r *http.Request, providerType string) (*OAuth2ProviderData, error) {
	query := r.URL.Query()

	codeQuery := query.Get("code")
	if codeQuery == "" {
		return nil, BadRequestError(ErrorCodeBadOAuth2Callback, "OAuth state parameter missing")
	}

	p, err := a.OAuthProvider(r.Context(), providerType)
	if err != nil {
		return nil, BadRequestError(ErrorCodeProviderOAuth2Unsupported, "unsupported provider: %+v", err)
	}

	token, err := p.GetOAuthToken(codeQuery)
	if err != nil {
		return nil, InternalServerError("unable to exchange code %s: %+v", codeQuery, token)
	}

	data, err := p.GetUserData(r.Context(), token)
	if err != nil {
		return nil, InternalServerError("error acquiring user data from external provider %s: %+v", providerType, err)
	}

	return &OAuth2ProviderData{
		userData:     data,
		accessToken:  token.AccessToken,
		refreshToken: token.RefreshToken,
		code:         codeQuery,
	}, nil
}

// OAuthProvider returns the corresponding oauth provider as an OAuthProvider interface
func (a *SurgeAPI) OAuthProvider(ctx context.Context, name string) (provider.OAuth2Provider, error) {
	providerCandidate, err := provider.Provider(ctx, a.config, name, "")
	if err != nil {
		return nil, err
	}

	switch p := providerCandidate.(type) {
	case provider.OAuth2Provider:
		return p, nil
	default:
		return nil, fmt.Errorf("provider %v cannot be used for OAuth", name)
	}
}
