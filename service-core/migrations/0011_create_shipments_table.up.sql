CREATE TYPE shipment_status_enum AS ENUM (
    'pending',
    'processing',
    'shipped',
    'in_transit',
    'delivered',
    'failed',
    'returned',
    'cancelled'
);

CREATE TABLE shipments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    transaction_id UUID NOT NULL UNIQUE,

    status shipment_status_enum NOT NULL DEFAULT 'pending',

    courier_name TEXT NOT NULL,
    tracking_number TEXT UNIQUE,

    shipping_cost NUMERIC(15, 2) NOT NULL,

    shipped_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT fk_transaction
        FOREIGN KEY(transaction_id)
        REFERENCES transactions(id)
        ON DELETE CASCADE
);