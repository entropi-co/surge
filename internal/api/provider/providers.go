package provider

import (
	"context"
	"fmt"
	"strings"
	"surge/internal/conf"
)

// Provider returns a Provider interface for the given name.
func Provider(ctx context.Context, config *conf.SurgeConfigurations, name string, scopes string) (OAuth2Provider, error) {
	name = strings.ToLower(name)

	switch name {
	case "google":
		return NewGoogleProvider(ctx, config.External.Google, scopes)
	default:
		return nil, fmt.Errorf("provider %s could not be found", name)
	}
}
