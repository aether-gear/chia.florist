CREATE TYPE payment_status
    AS ENUM (
        'pending',
        'paid',
        'failed',
        'expired',
        'cancelled',
        'refunded'
    );

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    order_id UUID NOT NULL,
    method_id UUID NOT NULL,
    payment_account_id UUID,

    provider TEXT NOT NULL,

    provider_payment_id TEXT,
    provider_order_id TEXT,

    amount BIGINT NOT NULL,

    status payment_status NOT NULL DEFAULT 'pending',

    expires_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    paid_at TIMESTAMP,

    CONSTRAINT fk_payment_order_id
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_payment_method_id
        FOREIGN KEY (method_id)
        REFERENCES payment_methods(id),

    CONSTRAINT fk_payment_account_id
        FOREIGN KEY (payment_account_id)
        REFERENCES payment_accounts(id)
);

CREATE INDEX idx_payments_order_id
    ON payments(order_id);

CREATE INDEX idx_payments_status
    ON payments(status);

CREATE INDEX idx_payments_provider_payment_id
    ON payments(provider_payment_id);
