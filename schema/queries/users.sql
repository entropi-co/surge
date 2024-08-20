-- name: CreateUser :one
INSERT INTO auth.users(id, email, username, encrypted_password, created_at, updated_at)
values ($1, $2, $3, $4, now(), now())
RETURNING *;