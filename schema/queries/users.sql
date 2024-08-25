-- name: CreateUser :one
insert into auth.users(email, username, encrypted_password, meta_avatar, meta_first_name, meta_last_name,
                       meta_birthdate,
                       meta_extra, created_at, updated_at)
values ($1, $2, $3, $4, $5, $6, $7, $8, now(), now())
returning *;

-- name: GetUser :one
select *
from auth.users
where id = $1;

-- name: GetUserByEmail :one
select *
from auth.users
where email = $1;

-- name: GetUserByUsername :one
select *
from auth.users
where username = $1;

-- name: UpdateUser :one
update auth.users
set email              = coalesce(sqlc.narg('email'), email),
    username           = coalesce(sqlc.narg('username'), username),
    encrypted_password = coalesce(sqlc.narg('encrypted_password'), encrypted_password),

    created_at         = coalesce(sqlc.narg('created_at'), created_at),
    updated_at         = now(),
    last_sign_in       = coalesce(sqlc.narg('last_sign_in'), last_sign_in)
where id = $1
returning *;

-- name: UpdateUserMetadata :one
update auth.users
set meta_avatar     = coalesce(sqlc.narg('meta_avatar'), meta_avatar),
    meta_first_name = coalesce(sqlc.narg('meta_first_name'), meta_first_name),
    meta_last_name  = coalesce(sqlc.narg('meta_last_name'), meta_first_name),
    meta_birthdate  = coalesce(sqlc.narg('meta_birthdate'), meta_first_name),
    meta_extra      = coalesce(sqlc.narg('meta_extra'), meta_extra)
where id = $1
returning *;


-- name: UpdateUserLastSignIn :one
update auth.users
set updated_at   = now(),
    last_sign_in = now()
where id = $1
returning *;