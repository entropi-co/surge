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