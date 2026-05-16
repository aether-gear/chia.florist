CREATE TYPE account_status
    AS ENUM ('pending', 'active', 'suspended', 'locked');

CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL UNIQUE,

    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,

    status account_status DEFAULT 'pending',

    last_login_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_email_per_account
    ON accounts(email)
    WHERE deleted_at IS NULL;