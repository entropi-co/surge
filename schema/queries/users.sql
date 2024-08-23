-- name: CreateUser :one
INSERT INTO auth.users(email, username, encrypted_password, created_at, updated_at)
values ($1, $2, $3, now(), now())
RETURNING *;

-- name: GetUser :one
SELECT *
from auth.users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT *
from auth.users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT *
from auth.users
WHERE username = $1;

-- name: UpdateUser :one
UPDATE auth.users
SET email              = coalesce(sqlc.narg('email'), email),
    username           = coalesce(sqlc.narg('username'), username),
    encrypted_password = coalesce(sqlc.narg('encrypted_password'), encrypted_password),
    created_at         = coalesce(sqlc.narg('created_at'), created_at),
    updated_at         = now(),
    last_sign_in       = coalesce(sqlc.narg('last_sign_in'), last_sign_in)
WHERE id = $1
RETURNING *;

-- name: UpdateUserLastSignIn :one
update auth.users
set updated_at   = now(),
    last_sign_in = now()
where id = $1
returning *;