package conf

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type JwkMap map[string]JwkPair

type JwkPair struct {
	PublicKey  jwk.Key `json:"public_key"`
	PrivateKey jwk.Key `json:"private_key"`
}

// Decode implements the Decoder interface
// From supabase/auth
func (j *JwkMap) Decode(value string) error {
	data := make([]json.RawMessage, 0)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	config := JwkMap{}
	for _, key := range data {
		privJwk, err := jwk.ParseKey(key)
		if err != nil {
			return err
		}
		pubJwk, err := jwk.PublicKeyOf(privJwk)
		if err != nil {
			return err
		}

		// all public keys should have the use claim set to 'sig
		if err := pubJwk.Set(jwk.KeyUsageKey, "sig"); err != nil {
			return err
		}

		// all public keys should only have 'verify' set as the key_ops
		if err := pubJwk.Set(jwk.KeyOpsKey, jwk.KeyOperationList{jwk.KeyOpVerify}); err != nil {
			return err
		}

		config[pubJwk.KeyID()] = JwkPair{
			PublicKey:  pubJwk,
			PrivateKey: privJwk,
		}
	}
	*j = config
	return nil
}

// Validate
// from supabase/auth
func (j *JwkMap) Validate() error {
	// Validate performs _minimal_ checks if the data stored in the key are valid.
	// By minimal, we mean that it does not check if the key is valid for use in
	// cryptographic operations. For example, it does not check if an RSA key's
	// `e` field is a valid exponent, or if the `n` field is a valid modulus.
	// Instead, it checks for things such as the _presence_ of some required fields,
	// or if certain keys' values are of particular length.
	//
	// Note that depending on the underlying key type, use of this method requires
	// that multiple fields in the key are properly populated. For example, an EC
	// key's "x", "y" fields cannot be validated unless the "crv" field is populated first.
	var signingKeys []jwk.Key
	for _, key := range *j {
		if err := key.PrivateKey.Validate(); err != nil {
			return err
		}
		// symmetric keys don't have public keys
		if key.PublicKey != nil {
			if err := key.PublicKey.Validate(); err != nil {
				return err
			}
		}

		for _, op := range key.PrivateKey.KeyOps() {
			if op == jwk.KeyOpSign {
				signingKeys = append(signingKeys, key.PrivateKey)
				break
			}
		}
	}

	switch {
	case len(signingKeys) == 0:
		return fmt.Errorf("no signing key detected")
	case len(signingKeys) > 1:
		return fmt.Errorf("multiple signing keys detected, only 1 signing key is supported")
	}

	return nil
}

func (c SurgeJWTConfigurations) GetSigningJwk() (jwk.Key, error) {
	for _, key := range c.Keys {
		for _, op := range key.PrivateKey.KeyOps() {
			// the private JWK with key_ops "sign" should be used as the signing key
			if op == jwk.KeyOpSign {
				return key.PrivateKey, nil
			}
		}
	}
	return nil, fmt.Errorf("no signing key found")
}

// GetJwkCompatibleAlgorithm converts jwx/v2/jwk key algorithm to jwt/v5 algorithm type
func GetJwkCompatibleAlgorithm(key jwk.Key) jwt.SigningMethod {
	if key == nil {
		return jwt.SigningMethodHS256
	}

	switch (key).Algorithm().String() {
	case "RS256":
		return jwt.SigningMethodRS256
	case "RS512":
		return jwt.SigningMethodRS512
	case "ES256":
		return jwt.SigningMethodES256
	case "ES512":
		return jwt.SigningMethodES512
	case "EdDSA":
		return jwt.SigningMethodEdDSA
	}

	// return HS256 to preserve existing behavior
	return jwt.SigningMethodHS256
}

func GetSigningKeyFromJwk(key jwk.Key) (any, error) {
	var raw any
	if err := key.Raw(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func GetPublicKeyByID(kid string, config *SurgeJWTConfigurations) (any, error) {
	if k, ok := config.Keys[kid]; ok {
		key, err := GetSigningKeyFromJwk(k.PublicKey)
		if err != nil {
			return nil, err
		}
		return key, nil
	}
	if kid == *config.KeyID {
		return []byte(config.Secret), nil
	}
	return nil, fmt.Errorf("invalid kid: %s", kid)
}
