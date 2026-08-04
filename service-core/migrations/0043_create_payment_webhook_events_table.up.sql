CREATE TABLE payment_webhook_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Idempotency key: one row per unique (order, status) transition.
    -- gateway_order_id mirrors the gateway's order id string (not the internal UUID)
    -- so no join is required when receiving the raw webhook.
    gateway_order_id TEXT NOT NULL,
    transaction_status TEXT NOT NULL,

    payload JSONB NOT NULL,

    -- received  → persisted, not yet processed
    -- processed → successfully applied to the payment domain
    -- failed    → processing failed; will be re-attempted on next delivery
    status TEXT NOT NULL DEFAULT 'received',

    -- Populated only when status = 'failed'
    error TEXT,

    received_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- Populated when status transitions to 'processed'
    processed_at TIMESTAMP,

    CONSTRAINT uq_payment_webhook_events_order_status
        UNIQUE (order_id, transaction_status)
);

CREATE INDEX idx_pwe_status ON payment_webhook_events(status);
CREATE INDEX idx_pwe_order_id ON payment_webhook_events(order_id);
