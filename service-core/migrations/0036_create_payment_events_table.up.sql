CREATE TABLE payment_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    payment_id UUID NOT NULL,

    event_name TEXT NOT NULL,

    payload JSONB NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_payment_events_payment_id
        FOREIGN KEY (payment_id)
        REFERENCES payments(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_payment_events_payment_id
    ON payment_events(payment_id);

CREATE INDEX idx_payment_events_event_name
    ON payment_events(event_name);