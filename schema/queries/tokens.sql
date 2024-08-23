-- name: CreateRefreshToken :one
insert into auth.refresh_tokens(user_id, token, revoked, created_at, updated_at)
values ($1, $2, $3, now(), now())
returning *;