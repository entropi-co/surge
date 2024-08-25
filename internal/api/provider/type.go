package provider

import (
	"context"
	"golang.org/x/oauth2"
)

// OAuth2Provider specifies additional methods needed for providers using OAuth
type OAuth2Provider interface {
	AuthCodeURL(string, ...oauth2.AuthCodeOption) string
	GetUserData(context.Context, *oauth2.Token) (*UserData, error)
	GetOAuthToken(string) (*oauth2.Token, error)
}

type UserClaims struct {
	// Reserved claims
	Issuer  string  `json:"iss,omitempty" structs:"iss,omitempty"`
	Subject string  `json:"sub,omitempty" structs:"sub,omitempty"`
	Aud     string  `json:"aud,omitempty" structs:"aud,omitempty"`
	Iat     float64 `json:"iat,omitempty" structs:"iat,omitempty"`
	Exp     float64 `json:"exp,omitempty" structs:"exp,omitempty"`

	// Default profile claims
	Name              string `json:"name,omitempty" structs:"name,omitempty"`
	FamilyName        string `json:"family_name,omitempty" structs:"family_name,omitempty"`
	GivenName         string `json:"given_name,omitempty" structs:"given_name,omitempty"`
	MiddleName        string `json:"middle_name,omitempty" structs:"middle_name,omitempty"`
	NickName          string `json:"nickname,omitempty" structs:"nickname,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty" structs:"preferred_username,omitempty"`
	Profile           string `json:"profile,omitempty" structs:"profile,omitempty"`
	Picture           string `json:"picture,omitempty" structs:"picture,omitempty"`
	Website           string `json:"website,omitempty" structs:"website,omitempty"`
	Gender            string `json:"gender,omitempty" structs:"gender,omitempty"`
	Birthdate         string `json:"birthdate,omitempty" structs:"birthdate,omitempty"`
	ZoneInfo          string `json:"zoneinfo,omitempty" structs:"zoneinfo,omitempty"`
	Locale            string `json:"locale,omitempty" structs:"locale,omitempty"`
	UpdatedAt         string `json:"updated_at,omitempty" structs:"updated_at,omitempty"`
	Email             string `json:"email,omitempty" structs:"email,omitempty"`
	EmailVerified     bool   `json:"email_verified,omitempty" structs:"email_verified"`
	Phone             string `json:"phone,omitempty" structs:"phone,omitempty"`
	PhoneVerified     bool   `json:"phone_verified,omitempty" structs:"phone_verified"`

	// Custom profile claims that are
	ExtraClaims map[string]interface{} `json:"custom_claims,omitempty" structs:"custom_claims,omitempty"`
}

type UserEmail struct {
	Email    string
	Verified bool
	Primary  bool
}

type UserData struct {
	Emails []UserEmail
	Claims *UserClaims
}

func (e UserEmail) String() string {
	return e.Email
}
