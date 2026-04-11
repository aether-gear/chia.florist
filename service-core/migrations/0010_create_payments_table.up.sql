CREATE TYPE payment_status_enum AS ENUM (
    'pending',
    'paid',
    'failed',
    'cancelled',
    'expired'
);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    transaction_id UUID NOT NULL,
    payment_method_id UUID NOT NULL,

    gateway_ref TEXT NOT NULL UNIQUE,

    status payment_status_enum NOT NULL DEFAULT 'pending',

    amount NUMERIC(15, 2) NOT NULL,

    paid_at TIMESTAMPTZ,
    expired_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_transaction
        FOREIGN KEY(transaction_id)
        REFERENCES transactions(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_payment_method
        FOREIGN KEY(payment_method_id)
        REFERENCES payment_methods(id)
        ON DELETE RESTRICT
);