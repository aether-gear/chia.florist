CREATE TABLE payment_channel_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    payment_id UUID UNIQUE NOT NULL,

    channel_type method_type NOT NULL,

    -- Human-readable label returned by the gateway
    -- e.g. "QRIS", "GoPay", "BCA Virtual Account"
    display_name TEXT NOT NULL,

    -- The actionable value: QR string, deep-link URL, or VA number.
    -- NULL when the gateway returns no instruction value.
    action_url TEXT,

    expires_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_payment_channel_data_payment_id
        FOREIGN KEY (payment_id)
        REFERENCES payments(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_payment_channel_data_payment_id
    ON payment_channel_data(payment_id);
