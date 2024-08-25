package api

type contextKey string

func (c contextKey) String() string {
	return "surge.api.context-key " + string(c)
}

const (
	contextExternalReferrerKey     = contextKey("external_referrer")
	contextTargetUserKey           = contextKey("target_user")
	contextExternalProviderTypeKey = contextKey("external_provider_type")
	contextSignatureKey            = contextKey("signature")
)
