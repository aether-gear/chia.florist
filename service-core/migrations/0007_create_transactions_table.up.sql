CREATE TYPE transaction_status_enum AS ENUM (
    'pending',
    'paid',
    'failed',
    'cancelled',
    'expired'
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL,
    address_id UUID NOT NULL,

    status transaction_status_enum NOT NULL DEFAULT 'pending',

    total_amount NUMERIC(15, 2) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    due_date TIMESTAMPTZ NOT NULL,
    paid_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    expired_at TIMESTAMPTZ,

    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_address
        FOREIGN KEY(address_id)
        REFERENCES addresses(id)
        ON DELETE CASCADE
);