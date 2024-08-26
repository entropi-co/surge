package api

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"regexp"
	"surge/internal/conf"
)

var authenticationBearerRegex = regexp.MustCompile(`^(?:B|b)earer (\S+$)`)

func (a *SurgeAPI) useAuthentication(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	token, err := a.getBearerAuthorizationHeader(r)
	config := a.config
	if err != nil {
		a.clearCookieTokens(config, w)
		return nil, err
	}

	ctx, err := a.parseJWTClaims(token, r)
	if err != nil {
		a.clearCookieTokens(config, w)
		return ctx, err
	}

	return ctx, err
}

func (a *SurgeAPI) getBearerAuthorizationHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	matches := authenticationBearerRegex.FindStringSubmatch(authHeader)
	if len(matches) != 2 {
		return "", UnauthorizedError(ErrorCodeNoAuthorization, "This endpoint requires a Bearer token")
	}

	return matches[1], nil
}

func (a *SurgeAPI) parseJWTClaims(bearer string, r *http.Request) (context.Context, error) {
	ctx := r.Context()
	config := a.config

	p := jwt.NewParser(jwt.WithValidMethods(config.JWT.ValidMethods))
	token, err := p.ParseWithClaims(bearer, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if kid, ok := token.Header["kid"]; ok {
			if kidStr, ok := kid.(string); ok {
				return conf.GetPublicKeyByID(kidStr, &config.JWT)
			}
		}
		if alg, ok := token.Header["alg"]; ok {
			if alg == jwt.SigningMethodHS256.Name {
				// preserve backward compatibility for cases where the kid is not set
				return []byte(config.JWT.Secret), nil
			}
		}
		return nil, fmt.Errorf("missing kid")
	})
	if err != nil {
		return nil, ForbiddenError(ErrorCodeBadJWT, "invalid JWT: unable to parse or verify signature, %v", err)
	}

	return context.WithValue(ctx, contextTokenKey, token), nil
}
