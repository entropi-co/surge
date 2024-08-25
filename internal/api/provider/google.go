package provider

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"strings"
	"surge/internal/conf"
)

type googleUser struct {
	ID            string `json:"id"`
	Subject       string `json:"sub"`
	Issuer        string `json:"iss"`
	Name          string `json:"name"`
	AvatarURL     string `json:"picture"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	EmailVerified bool   `json:"email_verified"`
	HostedDomain  string `json:"hd"`
}

type googleOAuth2Provider struct {
	*oauth2.Config

	oidc *oidc.Provider
}

const IssuerGoogle = "https://accounts.google.com"

var internalIssuerGoogle = IssuerGoogle

func NewGoogleProvider(ctx context.Context, config conf.SurgeProviderConfiguration, scopes string) (OAuth2Provider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	oauthScopes := []string{
		"email",
		"profile",
	}

	if scopes != "" {
		oauthScopes = append(oauthScopes, strings.Split(scopes, ",")...)
	}

	oidcProvider, err := oidc.NewProvider(ctx, internalIssuerGoogle)
	if err != nil {
		return nil, err
	}

	return &googleOAuth2Provider{
		Config: &oauth2.Config{
			ClientID:     config.ClientID[0],
			ClientSecret: config.ClientSecret,
			Endpoint:     oidcProvider.Endpoint(),
			Scopes:       oauthScopes,
			RedirectURL:  config.RedirectURI,
		},
		oidc: oidcProvider,
	}, nil
}

func (g googleOAuth2Provider) GetOAuthToken(code string) (*oauth2.Token, error) {
	return g.Exchange(context.Background(), code)
}

const UserInfoEndpointGoogle = "https://www.googleapis.com/userinfo/v2/me"

var internalUserInfoEndpointGoogle = UserInfoEndpointGoogle

func (g googleOAuth2Provider) GetUserData(ctx context.Context, tok *oauth2.Token) (*UserData, error) {
	if idToken := tok.Extra("id_token"); idToken != nil {
		_, data, err := ParseIDToken(ctx, g.oidc, &oidc.Config{
			ClientID: g.Config.ClientID,
		}, idToken.(string), ParseIDTokenOptions{
			AccessToken: tok.AccessToken,
		})
		if err != nil {
			return nil, err
		}

		return data, nil
	}

	// This whole section offers legacy support in case the Google OAuth2
	// flow does not return an ID Token for the user, which appears to
	// always be the case.
	logrus.Info("Using Google OAuth2 user info endpoint, an ID token was not returned by Google")

	var u googleUser
	if err := makeRequest(ctx, tok, g.Config, internalUserInfoEndpointGoogle, &u); err != nil {
		return nil, err
	}

	var data UserData

	if u.Email != "" {
		data.Emails = append(data.Emails, UserEmail{
			Email:    u.Email,
			Verified: u.IsEmailVerified(),
			Primary:  true,
		})
	}

	data.Claims = &UserClaims{
		Issuer:        internalUserInfoEndpointGoogle,
		Subject:       u.ID,
		Name:          u.Name,
		Picture:       u.AvatarURL,
		Email:         u.Email,
		EmailVerified: u.IsEmailVerified(),
	}

	return &data, nil
}

func (u googleUser) IsEmailVerified() bool {
	return u.VerifiedEmail || u.EmailVerified
}
