package provider

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
)

type ParseIDTokenOptions struct {
	SkipAccessTokenCheck bool
	AccessToken          string
}

func ParseIDToken(ctx context.Context, provider *oidc.Provider, config *oidc.Config, idToken string, options ParseIDTokenOptions) (*oidc.IDToken, *UserData, error) {
	if config == nil {
		config = &oidc.Config{
			// aud claim check to be performed by other flows
			SkipClientIDCheck: true,
		}
	}

	verifier := provider.VerifierContext(ctx, config)

	token, err := verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, nil, err
	}

	var data *UserData

	switch token.Issuer {
	case IssuerGoogle:
		token, data, err = parseGoogleIDToken(token)
	default:
		token, data, err = parseGenericIDToken(token)
	}

	if err != nil {
		return nil, nil, err
	}

	if !options.SkipAccessTokenCheck && token.AccessTokenHash != "" {
		if err := token.VerifyAccessToken(options.AccessToken); err != nil {
			return nil, nil, err
		}
	}

	return token, data, nil
}

//
// Start ID Token Parsers
//

func parseGenericIDToken(token *oidc.IDToken) (*oidc.IDToken, *UserData, error) {
	var data UserData

	if err := token.Claims(&data.Claims); err != nil {
		return nil, nil, err
	}

	if data.Claims.Email != "" {
		data.Emails = append(data.Emails, UserEmail{
			Email:    data.Claims.Email,
			Verified: data.Claims.EmailVerified,
			Primary:  true,
		})
	}

	if len(data.Emails) <= 0 {
		return nil, nil, fmt.Errorf("provider: Generic OIDC ID token from issuer %q must contain an email address", token.Issuer)
	}

	return token, &data, nil
}

func parseGoogleIDToken(token *oidc.IDToken) (*oidc.IDToken, *UserData, error) {
	var claims googleUser
	if err := token.Claims(&claims); err != nil {
		return nil, nil, err
	}

	var data UserData

	if claims.Email != "" {
		data.Emails = append(data.Emails, UserEmail{
			Email:    claims.Email,
			Verified: claims.IsEmailVerified(),
			Primary:  true,
		})
	}

	data.Claims = &UserClaims{
		Issuer:  claims.Issuer,
		Subject: claims.Subject,
		Name:    claims.Name,
		Picture: claims.AvatarURL,
	}

	if claims.HostedDomain != "" {
		data.Claims.ExtraClaims = map[string]any{
			"hd": claims.HostedDomain,
		}
	}

	return token, &data, nil
}
