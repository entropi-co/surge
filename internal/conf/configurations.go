package conf

import (
	"errors"
	"github.com/gobwas/glob"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"surge/internal/utilities"
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

	AutoLinkSameEmail bool `default:"true" split_words:"true"`
}

type SurgeJWTConfigurations struct {
	ExpiresAfter int `required:"true" split_words:"true"`

	Secret string `required:"true"`
	Keys   JwkMap
	KeyID  *string `split_words:"true"`
}

type SurgeLoggingConfigurations struct {
	EnableDebug   bool `split_words:"true"`
	EnableRequest bool `split_words:"true"`
}

type SurgeConfigurations struct {
	Auth     SurgeAuthenticateConfigurations `split_words:"true"`
	JWT      SurgeJWTConfigurations          `split_words:"true"`
	Database SurgeDatabaseConfigurations     `required:"true"`
	External SurgeExternalConfigurations

	ServiceURL      string               `required:"true" split_words:"true"`
	URIAllowListMap map[string]glob.Glob `split_words:"true"`
	URIAllowList    []string             `json:"uri_allow_list" split_words:"true"`

	Logging SurgeLoggingConfigurations `split_words:"true"`
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

	err := config.ApplyDefaults()
	if err != nil {
		return nil, err
	}

	return config, nil
}

// ApplyDefaults apply defaults to SurgeConfigurations
func (c *SurgeConfigurations) ApplyDefaults() error {
	if c.JWT.Keys == nil || len(c.JWT.Keys) == 0 {
		// Keys are not provided, default to secret
		if c.JWT.Secret == "" {
			return errors.New("JWT_SECRET is required if JWT_KEYS is not provided")
		}

		privateKey, err := jwk.FromRaw([]byte(c.JWT.Secret))
		if err != nil {
			return err
		}
		if c.JWT.KeyID != nil && *c.JWT.KeyID != "" {
			// Override key id
			if err := privateKey.Set(jwk.KeyIDKey, c.JWT.KeyID); err != nil {
				return err
			}
		}

		// Default algorithm to HS256
		if privateKey.Algorithm().String() == "" {
			if err := privateKey.Set(jwk.AlgorithmKey, jwt.SigningMethodHS256.Name); err != nil {
				return err
			}
		}
		if err := privateKey.Set(jwk.KeyUsageKey, "sig"); err != nil {
			return err
		}
		if len(privateKey.KeyOps()) == 0 {
			if err := privateKey.Set(jwk.KeyOpsKey, jwk.KeyOperationList{jwk.KeyOpSign, jwk.KeyOpVerify}); err != nil {
				return err
			}
		}

		// Create a public key from the private key
		pubKey, err := privateKey.PublicKey()
		if err != nil {
			return err
		}

		// Assign default JWKs
		c.JWT.Keys = make(JwkMap)
		c.JWT.Keys[*utilities.OrDefaultFn(c.JWT.KeyID, func() *string { return new(string) })] = JwkPair{
			PublicKey:  pubKey,
			PrivateKey: privateKey,
		}
	}

	if c.URIAllowList == nil {
		c.URIAllowList = []string{}
	}

	if c.URIAllowList != nil {
		c.URIAllowListMap = make(map[string]glob.Glob)
		for _, uri := range c.URIAllowList {
			g := glob.MustCompile(uri, '.', '/')
			c.URIAllowListMap[uri] = g
		}
	}
	return nil
}
