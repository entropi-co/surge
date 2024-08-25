// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: identities.sql

package schema

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

const createIdentityWithUser = `-- name: CreateIdentityWithUser :one
INSERT INTO auth.identities(user_id, provider, provider_id, provider_data, data, created_at, updated_at, last_sign_in)
values ($1, $2, $3, $4, '{}', now(), now(), null)
RETURNING id, user_id, data, provider, provider_id, provider_data, created_at, updated_at, last_sign_in
`

type CreateIdentityWithUserParams struct {
	UserID       uuid.UUID
	Provider     string
	ProviderID   string
	ProviderData json.RawMessage
}

func (q *Queries) CreateIdentityWithUser(ctx context.Context, arg CreateIdentityWithUserParams) (*AuthIdentity, error) {
	row := q.db.QueryRowContext(ctx, createIdentityWithUser,
		arg.UserID,
		arg.Provider,
		arg.ProviderID,
		arg.ProviderData,
	)
	var i AuthIdentity
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Data,
		&i.Provider,
		&i.ProviderID,
		&i.ProviderData,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastSignIn,
	)
	return &i, err
}

const getIdentitiesByUser = `-- name: GetIdentitiesByUser :many
SELECT id, user_id, data, provider, provider_id, provider_data, created_at, updated_at, last_sign_in
from auth.identities
WHERE user_id = $1
`

func (q *Queries) GetIdentitiesByUser(ctx context.Context, userID uuid.UUID) ([]*AuthIdentity, error) {
	rows, err := q.db.QueryContext(ctx, getIdentitiesByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*AuthIdentity
	for rows.Next() {
		var i AuthIdentity
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Data,
			&i.Provider,
			&i.ProviderID,
			&i.ProviderData,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.LastSignIn,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getIdentity = `-- name: GetIdentity :one
SELECT id, user_id, data, provider, provider_id, provider_data, created_at, updated_at, last_sign_in
from auth.identities
where provider = $1
  and provider_id = $2
`

type GetIdentityParams struct {
	Provider   string
	ProviderID string
}

func (q *Queries) GetIdentity(ctx context.Context, arg GetIdentityParams) (*AuthIdentity, error) {
	row := q.db.QueryRowContext(ctx, getIdentity, arg.Provider, arg.ProviderID)
	var i AuthIdentity
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Data,
		&i.Provider,
		&i.ProviderID,
		&i.ProviderData,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastSignIn,
	)
	return &i, err
}

const getIdentityById = `-- name: GetIdentityById :one
SELECT id, user_id, data, provider, provider_id, provider_data, created_at, updated_at, last_sign_in
from auth.identities
WHERE id = $1
`

func (q *Queries) GetIdentityById(ctx context.Context, id uuid.UUID) (*AuthIdentity, error) {
	row := q.db.QueryRowContext(ctx, getIdentityById, id)
	var i AuthIdentity
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Data,
		&i.Provider,
		&i.ProviderID,
		&i.ProviderData,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastSignIn,
	)
	return &i, err
}

const updateIdentity = `-- name: UpdateIdentity :one
UPDATE auth.identities
SET created_at    = coalesce($2, created_at),
    updated_at    = now(),
    provider_data = coalesce($3, provider_data),
    provider_id   = coalesce($4, provider_id),
    data          = coalesce($5, data),
    last_sign_in  = coalesce($6, last_sign_in)
WHERE id = $1
RETURNING id, user_id, data, provider, provider_id, provider_data, created_at, updated_at, last_sign_in
`

type UpdateIdentityParams struct {
	ID           uuid.UUID
	CreatedAt    sql.NullTime
	ProviderData pqtype.NullRawMessage
	ProviderID   sql.NullString
	Data         pqtype.NullRawMessage
	LastSignIn   sql.NullTime
}

func (q *Queries) UpdateIdentity(ctx context.Context, arg UpdateIdentityParams) (*AuthIdentity, error) {
	row := q.db.QueryRowContext(ctx, updateIdentity,
		arg.ID,
		arg.CreatedAt,
		arg.ProviderData,
		arg.ProviderID,
		arg.Data,
		arg.LastSignIn,
	)
	var i AuthIdentity
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Data,
		&i.Provider,
		&i.ProviderID,
		&i.ProviderData,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastSignIn,
	)
	return &i, err
}

const updateIdentityLastSignIn = `-- name: UpdateIdentityLastSignIn :one
UPDATE auth.identities
Set updated_at   = now(),
    last_sign_in = now()
WHERE id = $1
RETURNING id, user_id, data, provider, provider_id, provider_data, created_at, updated_at, last_sign_in
`

func (q *Queries) UpdateIdentityLastSignIn(ctx context.Context, id uuid.UUID) (*AuthIdentity, error) {
	row := q.db.QueryRowContext(ctx, updateIdentityLastSignIn, id)
	var i AuthIdentity
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Data,
		&i.Provider,
		&i.ProviderID,
		&i.ProviderData,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastSignIn,
	)
	return &i, err
}