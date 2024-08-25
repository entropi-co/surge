package conf

import "errors"

var (
	ErrProviderDisabled = errors.New("provider is disabled")
	ErrNoClientID       = errors.New("missing client id")
	ErrNoClientSecret   = errors.New("missing client secret")
	ErrNoRedirectUri    = errors.New("missing redirect uri")
)

type SurgeProviderConfiguration struct {
	ClientID     []string `json:"client_id" split_words:"true"`
	ClientSecret string   `json:"client_secret"  split_words:"true"`
	RedirectURI  string   `json:"redirect_uri" split_words:"true"`
	URL          string   `json:"url" split_words:"true"`
	ApiURL       string   `json:"api_url" split_words:"true"`
	Enabled      bool     `json:"enabled" default:"true"`
}

func (c SurgeProviderConfiguration) Validate() error {
	if !c.Enabled {
		return ErrProviderDisabled
	}
	if len(c.ClientID) == 0 {
		return ErrNoClientID
	}
	if c.ClientSecret == "" {
		return ErrNoClientSecret
	}
	if c.RedirectURI == "" {
		return ErrNoRedirectUri
	}
	return nil
}

type SurgeExternalConfigurations struct {
	Google SurgeProviderConfiguration `json:"google"`
}
