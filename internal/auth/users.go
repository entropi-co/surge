package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"surge/internal/api/provider"
	"surge/internal/conf"
	"surge/internal/schema"
	"surge/internal/storage"
	"time"
)

type UserMetadata struct {
	Avatar    *string
	FirstName *string
	LastName  *string
	Birthdate *time.Time
	Extra     map[string]interface{}
}

type CreateUserOptions struct {
	Email    *string `validate:"email,lte=255"`
	Username *string `validate:"gte=3,lte=20"`
	Phone    *string `validate:"e164"`
	Password *string `validate:"required,gte=8,lte=255"`
	Metadata UserMetadata
}

func (o CreateUserOptions) validate() error {
	if o.Email == nil && o.Username == nil {
		return ErrMissingField
	}

	validate := validator.New()

	var fieldsToExclude []string
	if o.Email == nil {
		fieldsToExclude = append(fieldsToExclude, "Email")
	}
	if o.Phone == nil {
		fieldsToExclude = append(fieldsToExclude, "Phone")
	}

	return validate.StructExcept(o, fieldsToExclude...)
}

type CreateUserAndIdentityOptions struct {
	Provider          string
	ProviderAccountID string
	ProviderData      provider.UserData
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
	if config.Auth.CredentialsRequirePhone && o.Phone == nil {
		return ErrRequiredPhone
	}

	return nil
}

func CreateUser(queries *schema.Queries, ctx context.Context, config *conf.SurgeConfigurations, options CreateUserOptions) (*schema.AuthUser, error) {
	if err := options.validateWithConfig(config); err != nil {
		return nil, err
	}

	if _, err := queries.GetUserByEmail(ctx, *options.Email); options.Email != nil {
		if err == nil {
			return nil, ErrDuplicateEmail
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDatabaseJob
		}
	}

	if _, err := queries.GetUserByUsername(ctx, *options.Username); options.Email != nil {
		if err == nil {
			return nil, ErrDuplicateUsername
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDuplicatePhone
		}
	}

	result, err := queries.CreateUser(ctx, schema.CreateUserParams{
		Phone:             storage.NewNullableString(options.Phone),
		Email:             storage.NewNullableString(options.Email),
		Username:          storage.NewNullableString(options.Username),
		EncryptedPassword: storage.NewNullableString(options.Password),

		MetaAvatar:    storage.NewNullableString(options.Metadata.Avatar),
		MetaFirstName: storage.NewNullableString(options.Metadata.FirstName),
		MetaLastName:  storage.NewNullableString(options.Metadata.LastName),
		MetaBirthdate: storage.NewNullableTime(options.Metadata.Birthdate),
		MetaExtra:     json.RawMessage("{}"),
	})
	if err != nil {
		logrus.WithError(err).Error(ErrDatabaseJob)
		return nil, ErrDatabaseJob
	}

	return result, err
}

func AuthenticateUser(user *schema.AuthUser, password string) bool {
	if !user.EncryptedPassword.Valid {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword.String), []byte(password)) == nil
}

func CreateUserAndIdentity(queries *schema.Queries, ctx context.Context, options CreateUserAndIdentityOptions) (*schema.AuthUser, *schema.AuthIdentity, error) {
	firstName := options.ProviderData.Claims.GivenName
	lastName := options.ProviderData.Claims.FamilyName
	avatarUrl := options.ProviderData.Claims.Picture

	user, err := queries.CreateUser(ctx, schema.CreateUserParams{
		MetaFirstName: sql.NullString{String: firstName, Valid: firstName != ""},
		MetaLastName:  sql.NullString{String: lastName, Valid: lastName != ""},
		MetaAvatar:    sql.NullString{String: avatarUrl, Valid: avatarUrl != ""},
	})
	if err != nil {
		return nil, nil, err
	}

	marshalledJson, err := json.Marshal(options.ProviderData)
	if err != nil {
		return nil, nil, err
	}

	identity, err := queries.CreateIdentityWithUser(ctx, schema.CreateIdentityWithUserParams{
		UserID:       user.ID,
		Provider:     options.Provider,
		ProviderID:   options.ProviderAccountID,
		ProviderData: marshalledJson,
	})
	if err != nil {
		return nil, nil, err
	}

	return user, identity, nil
}
