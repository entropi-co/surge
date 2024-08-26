package api

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

func (c contextKey) String() string {
	return "surge.api.context-key " + string(c)
}

const (
	contextExternalReferrerKey     = contextKey("external_referrer")
	contextTargetUserKey           = contextKey("target_user")
	contextExternalProviderTypeKey = contextKey("external_provider_type")
	contextSignatureKey            = contextKey("signature")
	contextTokenKey                = contextKey("token")
)

// getToken reads the JWT token from the context.
func getToken(ctx context.Context) *jwt.Token {
	obj := ctx.Value(contextTokenKey)
	if obj == nil {
		return nil
	}

	return obj.(*jwt.Token)
}

func getClaims(ctx context.Context) *AccessTokenClaims {
	token := getToken(ctx)
	if token == nil {
		return nil
	}
	return token.Claims.(*AccessTokenClaims)
}
