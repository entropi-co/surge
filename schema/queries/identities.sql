-- name: CreateIdentityWithUser :one
INSERT INTO auth.identities(user_id, provider, provider_id, provider_data, data, created_at, updated_at, last_sign_in)
values ($1, $2, $3, $4, '{}', now(), now(), null)
RETURNING *;

-- name: GetIdentity :one
SELECT *
from auth.identities
where provider = $1
  and provider_id = $2;

-- name: GetIdentityById :one
SELECT *
from auth.identities
WHERE id = $1;

-- name: GetIdentitiesByUser :many
SELECT *
from auth.identities
WHERE user_id = $1;

-- name: UpdateIdentity :one
UPDATE auth.identities
SET created_at    = coalesce(sqlc.narg('created_at'), created_at),
    updated_at    = now(),
    provider_data = coalesce(sqlc.narg('provider_data'), provider_data),
    provider_id   = coalesce(sqlc.narg('provider_id'), provider_id),
    data          = coalesce(sqlc.narg('data'), data),
    last_sign_in  = coalesce(sqlc.narg('last_sign_in'), last_sign_in)
WHERE id = $1
RETURNING *;

-- name: UpdateIdentityLastSignIn :one
UPDATE auth.identities
Set updated_at   = now(),
    last_sign_in = now()
WHERE id = $1
RETURNING *;