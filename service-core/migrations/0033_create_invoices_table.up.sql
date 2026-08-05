CREATE TYPE invoice_status
    AS ENUM (
        'issued',
        'void'
    );

CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    number VARCHAR(50) UNIQUE NOT NULL,
    order_id UUID UNIQUE NOT NULL,

    status invoice_status NOT NULL,

    subtotal BIGINT NOT NULL,
    shipping_fee BIGINT NOT NULL,
    total BIGINT NOT NULL,

    issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_invoices_order_id
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE,

    CONSTRAINT invoices_subtotal_check
        CHECK (subtotal >= 0),

    CONSTRAINT invoices_shipping_fee_check
        CHECK (shipping_fee >= 0),

    CONSTRAINT invoices_total_check
        CHECK (total >= 0)
);
