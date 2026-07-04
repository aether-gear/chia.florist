CREATE TYPE shipment_status_enum AS ENUM (
    'created',
    'packed',
    'labelled',
    'picked_up',
    'in_transit',
    'out_for_delivery',
    'delivered',
    'failed',
    'returned',
    'cancelled'
);

CREATE TABLE shipments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    order_id UUID NOT NULL UNIQUE,

    status shipment_status_enum NOT NULL DEFAULT 'created',
    tracking_number TEXT UNIQUE,

    courier_name TEXT NOT NULL,
    service TEXT NOT NULL,

    shipping_cost NUMERIC(15, 2) NOT NULL,
    weight TEXT NOT NULL,
    origin_id TEXT NOT NULL,
    destination_id TEXT NOT NULL,

    shipped_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT fk_order
        FOREIGN KEY(order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE
);
