package conf

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type SurgeDatabaseConfigurations struct {
	Url string `required:"true"`
}

type SurgeAuthenticateConfigurations struct {
	CredentialsRequireEmail    bool `default:"false" split_words:"true"`
	CredentialsRequirePhone    bool `default:"false" split_words:"true"`
	CredentialsRequireUsername bool `default:"false" split_words:"true"`

	DisableEmailAuth    bool `default:"false" split_words:"true"`
	DisableUsernameAuth bool `default:"false" split_words:"true"`
	DisablePhoneAuth    bool `default:"false" split_words:"true"`
}

type SurgeJWTConfigurations struct {
	ExpiresAfter int `required:"true" split_words:"true"`

	Keys JwkMap `json:"keys"`
}

type SurgeConfigurations struct {
	Auth     SurgeAuthenticateConfigurations `split_words:"true"`
	JWT      SurgeJWTConfigurations          `split_words:"true"`
	Database SurgeDatabaseConfigurations     `required:"true"`
}

func LoadFromEnvironments() (*SurgeConfigurations, error) {
	// Load .env
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	config := new(SurgeConfigurations)
	if err := envconfig.Process("surge", config); err != nil {
		return nil, err
	}

	return config, nil
}
