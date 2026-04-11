CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name TEXT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    phone TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    last_login_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX unique_username_per_users ON users(username)
WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX unique_email_per_users ON users(email)
WHERE deleted_at IS NULL;