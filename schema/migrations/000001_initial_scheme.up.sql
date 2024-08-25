create extension pgcrypto;

create schema if not exists auth;

create table if not exists auth.users
(
    id                 uuid                     not null unique default gen_random_uuid(),

    email              varchar(255)             null unique,
    username           varchar(255)             null unique,
    encrypted_password varchar(255)             null,

    meta_avatar        text                     null,
    meta_first_name    varchar(50)              null,
    meta_last_name     varchar(50)              null,
    meta_birthdate     date                     null,
    meta_extra         jsonb                    not null        default '{}',

    created_at         timestamp with time zone not null,
    updated_at         timestamp with time zone not null,

    last_sign_in       timestamp with time zone null            default null
);
create index if not exists users_id_email_index ON auth.users using brin (id);
create index if not exists users_id_username_index ON auth.users using brin (username);

create table if not exists auth.refresh_tokens
(
    id         bigserial                not null,
    user_id    uuid                     null,

    token      varchar(255)             null,
    revoked    bool                     not null default false,

    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,

    constraint refresh_tokens_pkey primary key (id)
);
create index if not exists refresh_tokens_id_index on auth.refresh_tokens using brin (id);
create index if not exists refresh_tokens_id_user_id_index on auth.refresh_tokens using brin (id, user_id);
create index if not exists refresh_tokens_token_index on auth.refresh_tokens using brin (token);

create table if not exists auth.identities
(
    id            uuid                     not null unique default gen_random_uuid(),
    user_id       uuid                     not null,

    data          jsonb                    not null,

    provider      text                     not null,
    provider_id   text                     not null,
    provider_data jsonb                    not null,

    created_at    timestamp with time zone not null,
    updated_at    timestamp with time zone not null,
    last_sign_in  timestamp with time zone null            default null,

    constraint identities_pkey primary key (id)
);
create index if not exists identities_id_index on auth.identities using brin (id);
create index if not exists identities_id_user_id_index on auth.identities using brin (id, user_id);