package api

import (
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"net/url"
	"strconv"
	"surge/internal/schema"
	"surge/internal/storage"
	"time"
)

// AccessTokenResponse represents an OAuth2 success response
type AccessTokenResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int           `json:"expires_in"`
	ExpiresAt    int64         `json:"expires_at"`
	User         *UserResponse `json:"user"`
}

// MakeRedirectUrl makes the token response to url so client can read token from query parameters
func (r *AccessTokenResponse) MakeRedirectUrl(redirectURL string, extraParams url.Values) string {
	extraParams.Set("access_token", r.AccessToken)
	extraParams.Set("expires_in", strconv.Itoa(r.ExpiresIn))
	extraParams.Set("expires_at", strconv.FormatInt(r.ExpiresAt, 10))
	extraParams.Set("refresh_token", r.RefreshToken)
	extraParams.Set("token_type", "bearer")

	return redirectURL + "#" + extraParams.Encode()
}

type UserResponse struct {
	ID uuid.UUID `json:"id"`

	Email    *string `json:"email"`
	Username *string `json:"username"`

	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	LastSignIn *time.Time `json:"last_sign_in"`
}

func NewUserResponse(user *schema.AuthUser) *UserResponse {
	return &UserResponse{
		ID:         user.ID,
		Email:      storage.NullStringToPointer(user.Email),
		Username:   storage.NullStringToPointer(user.Username),
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		LastSignIn: storage.NullTimeToPointer(user.LastSignIn),
	}
}

// JwksResponse is response type for /.well-known/jwks.json endpoint
type JwksResponse struct {
	Keys []jwk.Key `json:"keys"`
}
