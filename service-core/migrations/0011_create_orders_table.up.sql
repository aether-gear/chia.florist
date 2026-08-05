CREATE TYPE order_status
    AS ENUM (
        'pending',
        'confirmed',
        'processing',
        'shipped',
        'delivered',
        'cancelled',
        'expired'
    );

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    number VARCHAR(50) UNIQUE NOT NULL,

    customer_id UUID NOT NULL,
    address_id UUID NOT NULL,

    status order_status NOT NULL,

    subtotal BIGINT NOT NULL,
    shipping_fee BIGINT NOT NULL,
    total BIGINT NOT NULL,

    confirmed_at TIMESTAMP,
    handling_expires_at TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,

    CONSTRAINT fk_orders_customer_id
        FOREIGN KEY (customer_id)
        REFERENCES customers(id),

    CONSTRAINT fk_orders_address_id
        FOREIGN KEY (address_id)
        REFERENCES customer_addresses(id),

    CONSTRAINT orders_subtotal_check
        CHECK (subtotal >= 0),

    CONSTRAINT orders_shipping_fee_check
        CHECK (shipping_fee >= 0),

    CONSTRAINT orders_total_check
        CHECK (total >= 0),

    CONSTRAINT check_orders_handling_sla_timestamps
        CHECK (
            (status NOT IN ('confirmed', 'processing')) OR
            (confirmed_at IS NOT NULL AND handling_expires_at IS NOT NULL)
        )
);

CREATE INDEX idx_orders_status_handling_expires_at
    ON orders(status, handling_expires_at)
    WHERE status IN ('confirmed', 'processing');
