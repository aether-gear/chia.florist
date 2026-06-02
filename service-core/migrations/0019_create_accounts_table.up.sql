CREATE TYPE account_status
    AS ENUM ('pending', 'active', 'suspended', 'locked');

CREATE TYPE account_type
    AS ENUM ('customer', 'merchant');

CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL UNIQUE,

    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,

    status account_status DEFAULT 'pending',
    type account_type NOT NULL,

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