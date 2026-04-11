create TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    transaction_id UUID NOT NULL UNIQUE,
    address_id UUID NOT NULL UNIQUE,

    invoice_number TEXT NOT NULL UNIQUE,

    total_amount NUMERIC(15, 2) NOT NULL,

    issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    due_date TIMESTAMPTZ NOT NULL,
    paid_at TIMESTAMPTZ,

    CONSTRAINT fk_transaction
        FOREIGN KEY(transaction_id)
        REFERENCES transactions(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_address
        FOREIGN KEY(address_id)
        REFERENCES addresses(id)
        ON DELETE CASCADE
)