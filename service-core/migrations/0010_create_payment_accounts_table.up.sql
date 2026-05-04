CREATE TABLE payment_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    method_id UUID NOT NULL,

    account_name TEXT NOT NULL,
    account_number TEXT NOT NULL,
    phone_number TEXT,
    qr_string TEXT,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    current_load INTEGER NOT NULL DEFAULT 0,

    last_used_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_payment_accounts_method
        FOREIGN KEY (method_id)
        REFERENCES payment_methods(id)
        ON DELETE RESTRICT,

    CONSTRAINT payment_account_identifier_check
        CHECK (
            account_number IS NOT NULL OR
            phone_number IS NOT NULL OR
            qr_string IS NOT NULL
        ),

    CONSTRAINT unique_method_account
        UNIQUE (method_id, account_number)
);

CREATE INDEX idx_payment_accounts_active_load
    ON payment_accounts(method_id, current_load, last_used_at)
    WHERE deleted_at IS NULL AND is_active = TRUE;