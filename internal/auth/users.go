package auth

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"surge/internal/conf"
	"surge/internal/schema"
	"surge/internal/storage"
)

type CreateUserOptions struct {
	Email    *string `validate:"email,lte=255"`
	Username *string `validate:"gte=3,lte=20"`
	Password *string `validate:"required,gte=8,lte=255"`
}

func (o CreateUserOptions) validate() error {
	if o.Email == nil && o.Username == nil {
		return ErrMissingField
	}

	validate := validator.New()

	if o.Email == nil {
		return validate.StructExcept(o, "Email")
	} else {
		return validate.Struct(o)
	}
}

func (o CreateUserOptions) validateWithConfig(config *conf.SurgeConfigurations) error {
	if err := o.validate(); err != nil {
		return err
	}

	if config.Auth.CredentialsRequireEmail && o.Email == nil {
		return ErrRequiredEmail
	}
	if config.Auth.CredentialsRequireUsername && o.Username == nil {
		return ErrRequiredUsername
	}
	// TODO: Validate if phone field is required

	return nil
}

func CreateUser(queries *schema.Queries, ctx context.Context, config *conf.SurgeConfigurations, options CreateUserOptions) (*schema.AuthUser, error) {
	if err := options.validateWithConfig(config); err != nil {
		return nil, err
	}

	if _, err := queries.GetUserByEmail(ctx, storage.NewNullableString(options.Email)); options.Email != nil {
		if err == nil {
			return nil, ErrDuplicateEmail
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDatabaseJob
		}
	}

	if _, err := queries.GetUserByUsername(ctx, storage.NewNullableString(options.Username)); options.Email != nil {
		if err == nil {
			return nil, ErrDuplicateUsername
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDuplicatePhone
		}
	}

	result, err := queries.CreateUser(ctx, schema.CreateUserParams{
		Email:             storage.NewNullableString(options.Email),
		Username:          storage.NewNullableString(options.Username),
		EncryptedPassword: storage.NewNullableString(options.Password),
	})
	if err != nil {
		return nil, ErrDatabaseJob
	}

	return &result, err
}

func AuthenticateUser(user *schema.AuthUser, password string) bool {
	if !user.EncryptedPassword.Valid {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword.String), []byte(password)) == nil
}
