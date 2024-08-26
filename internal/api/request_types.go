package api

import "time"

type SignUpWithCredentialsRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Password *string `json:"password"`

	Metadata struct {
		Avatar    *string    `json:"avatar"`
		FirstName *string    `json:"first_name"`
		LastName  *string    `json:"last_name"`
		Birthdate *time.Time `json:"birthdate"`
	}
}
