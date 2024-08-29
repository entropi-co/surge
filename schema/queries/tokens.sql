-- name: CreateRefreshToken :one
insert into auth.refresh_tokens(user_id, token, revoked, created_at, updated_at)
values ($1, $2, $3, now(), now())
returning *;

-- name: ListRefreshTokenByUser :many
select *
from auth.refresh_tokens
where user_id = $1;

-- name: GetRefreshToken :one
select *
from auth.refresh_tokens
where token = sqlc.arg('token')::varchar;

-- name: RevokeRefreshToken :exec
update auth.refresh_tokens
set revoked = true
where id = $1;

-- name: RevokeRefreshTokensOfUser :exec
update auth.refresh_tokens
set revoked = true
where user_id = $1;