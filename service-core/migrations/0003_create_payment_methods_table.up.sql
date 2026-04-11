CREATE TABLE payment_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,

    provider TEXT,

    is_active BOOLEAN DEFAULT TRUE,

    fee_percentage NUMERIC(5,2),
    fee_flat NUMERIC(15,2),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);