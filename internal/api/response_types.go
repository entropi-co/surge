package api

import (
	"github.com/google/uuid"
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
		Username:   storage.NullStringToPointer(user.Email),
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		LastSignIn: storage.NullTimeToPointer(user.LastSignIn),
	}
}
