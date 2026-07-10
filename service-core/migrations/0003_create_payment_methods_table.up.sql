CREATE TYPE fee_type AS ENUM ('flat', 'percentage', 'mixed');

CREATE TYPE method_type AS ENUM ('bank_transfer', 'ewallet', 'qr_code');

CREATE TYPE method_code AS ENUM ('gopay', 'shopeepay', 'qris', 'bca_va', 'mandiri_bill');

CREATE TABLE payment_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name TEXT NOT NULL,
    code TEXT NOT NULL,
    provider TEXT NOT NULL,
    type method_type NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    description TEXT,

    fee_type fee_type NOT NULL,
    fee_amount NUMERIC(15,2) DEFAULT 0,
    fee_rate NUMERIC(5,4) DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT payment_methods_code_provider_type_key UNIQUE (code, provider, type),

    CONSTRAINT payment_methods_fee_check
        CHECK (
            (fee_type = 'flat' AND fee_rate = 0)
         OR (fee_type = 'percentage' AND fee_amount = 0)
         OR (fee_type = 'mixed')
        ),

    CONSTRAINT payment_methods_fee_amount_non_negative
        CHECK (fee_amount >= 0),

    CONSTRAINT payment_methods_fee_rate_range
        CHECK (fee_rate >= 0 AND fee_rate <= 1)
);
