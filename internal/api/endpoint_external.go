package api

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"strconv"
	"surge/internal/api/provider"
	"surge/internal/auth"
	"surge/internal/schema"
	"time"
)

type ExternalProviderClaims struct {
	jwt.RegisteredClaims
	Provider        string `json:"provider"`
	Referrer        string `json:"referrer,omitempty"`
	LinkingTargetID string `json:"linking_target_id,omitempty"`
}

func (a *SurgeAPI) EndpointExternal(w http.ResponseWriter, r *http.Request) error {
	targetUrl, err := a.GetExternalProviderUrl(w, r)
	if err != nil {
		return err
	}

	if parsed, _ := strconv.ParseBool(r.URL.Query().Get("no_redirect")); !parsed {
		http.Redirect(w, r, targetUrl, http.StatusFound)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(targetUrl))
	}

	return err
}

func (a *SurgeAPI) EndpointExternalCallback(w http.ResponseWriter, r *http.Request) error {
	a.redirectErrors(a.internalExternalProviderCallback, w, r)
	return nil
}

func (a *SurgeAPI) internalExternalProviderCallback(w http.ResponseWriter, r *http.Request) error {
	providerType := r.Context().Value(contextExternalProviderTypeKey).(string)
	//logger := logrus.WithContext(r.Context()).WithField("req", middleware.GetReqID(r.Context()))

	data, err := a.handleOAuthCallback(r)
	if err != nil {
		return err
	}

	userData := data.userData

	if len(userData.Emails) == 0 {
		return InternalServerError("failed to acquire user emails from provider %s: there was no email given", providerType)
	}

	// Reset
	userData.Claims.EmailVerified = false
	for _, email := range userData.Emails {
		userData.Claims.Email = email.Email
		userData.Claims.EmailVerified = email.Verified

		if email.Primary {
			break
		}
	}

	var user *schema.AuthUser
	var identity *schema.AuthIdentity
	identity, err = a.queries.GetIdentity(r.Context(), schema.GetIdentityParams{
		Provider:   providerType,
		ProviderID: userData.Claims.Subject,
	})
	if err != nil {
		if !errors.Is(sql.ErrNoRows, err) {
			return InternalServerError("database failed to find existing identity")
		}

		user, identity, err = auth.CreateUserAndIdentity(a.queries, r.Context(), auth.CreateUserAndIdentityOptions{
			Provider:          providerType,
			ProviderAccountID: userData.Claims.Subject,
			ProviderData:      *userData,
		})
		if err != nil {
			return InternalServerError("database failed to create new user and identity: %+v", err)
		}
	} else {
		user, err = a.queries.GetUser(r.Context(), identity.UserID)
		if err != nil {
			return InternalServerError("database failed to get user from existing identity")
		}
	}

	token, err := a.issueToken(r.Context(), user)
	if err != nil {
		return InternalServerError("failed to issue token")
	}

	redirectTo := a.getExternalRedirectURL(r)

	q := url.Values{}
	q.Set("provider_access_token", data.accessToken)
	if data.refreshToken != "" {
		q.Set("provider_refresh_token", data.refreshToken)
	}

	redirectTo = token.MakeRedirectUrl(redirectTo, q)

	logrus.Debugln(redirectTo)

	http.Redirect(w, r, redirectTo, http.StatusFound)
	return nil
}

func (a *SurgeAPI) GetExternalProviderUrl(w http.ResponseWriter, r *http.Request) (string, error) {
	query := r.URL.Query()

	providerType := query.Get("provider")
	scopes := query.Get("scopes")
	//codeChallenge := query.Get("code_challenge")
	//codeChallengeMethod := query.Get("code_challenge_method")

	p, err := provider.Provider(r.Context(), a.config, providerType, scopes)
	if err != nil {
		return "", BadRequestError(ErrorCodeInvalidProviderType, "unsupported provider: %+v", err)
	}

	claims := ExternalProviderClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(5 * time.Minute)},
		},
		Provider: providerType,
		Referrer: GetRequestReferrer(r, a.config),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(a.config.JWT.Secret))
	if err != nil {
		return "", InternalServerError("error creating state: %+v", err)
	}

	query.Del("scopes")
	query.Del("provider")
	authUrlParams := make([]oauth2.AuthCodeOption, 0)
	for key := range query {
		authUrlParams = append(authUrlParams, oauth2.SetAuthURLParam(key, query.Get("key")))
	}

	authUrl := p.AuthCodeURL(signedToken, authUrlParams...)

	return authUrl, nil
}

func (a *SurgeAPI) redirectErrors(handler surgeAPIHandler, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	errorID := middleware.GetReqID(ctx)
	err := handler(w, r)
	if err != nil {
		q := getErrorQueryString(err, errorID, logrus.New())
		http.Redirect(w, r, a.getExternalRedirectURL(r)+"#"+q.Encode(), http.StatusFound)
	}
}

func getErrorQueryString(err error, errorID string, log logrus.FieldLogger) *url.Values {
	q := url.Values{}
	var e *HTTPError
	switch {
	case errors.As(err, &e):
		q.Set("error", "server_error")

		if e.HTTPStatus >= http.StatusInternalServerError {
			e.ErrorID = errorID
			// this will get us the stack trace too
			log.WithError(e.Cause()).Error(e.Error())
		} else {
			log.WithError(e.Cause()).Info(e.Error())
		}
		q.Set("error_description", e.Message)
	default:
		q.Set("error", "server_error")
		q.Set("error_description", err.Error())
	}
	return &q
}

func (a *SurgeAPI) getExternalRedirectURL(r *http.Request) string {
	ctx := r.Context()
	config := a.config
	if er := ctx.Value(contextExternalReferrerKey).(string); er != "" {
		return er
	}
	return config.ServiceURL
}

func (a *SurgeAPI) loadExternalStateToContext(ctx context.Context, state string) (context.Context, error) {
	claims := ExternalProviderClaims{}
	p := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	_, err := p.ParseWithClaims(state, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.config.JWT.Secret), nil
	})
	if err != nil {
		return nil, BadRequestError(ErrorCodeBadOAuth2State, "OAuth callback with invalid state, %+v", err)
	}
	if claims.Provider == "" {
		return nil, BadRequestError(ErrorCodeBadOAuth2State, "OAuth callback with invalid state (missing provider)")
	}
	if claims.Referrer != "" {
		ctx = context.WithValue(ctx, contextExternalReferrerKey, claims.Referrer)
	}
	if claims.LinkingTargetID != "" {
		linkingTargetUserID, err := uuid.Parse(claims.LinkingTargetID)
		if err != nil {
			return nil, BadRequestError(ErrorCodeBadOAuth2State, "OAuth callback with invalid state (linking_target_id must be UUID)")
		}
		u, err := a.queries.GetUser(ctx, linkingTargetUserID)
		if err != nil {
			if errors.Is(sql.ErrNoRows, err) {
				return nil, UnprocessableEntityError(ErrorCodeUserNotFound, "Linking target user not found")
			}
			return nil, InternalServerError("Database error loading user: %+v", err)
		}
		ctx = context.WithValue(ctx, contextTargetUserKey, u)
	}
	ctx = context.WithValue(ctx, contextExternalProviderTypeKey, claims.Provider)
	return context.WithValue(ctx, contextSignatureKey, state), nil
}
