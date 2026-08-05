CREATE TABLE payment_instructions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_method_id UUID UNIQUE NOT NULL,

    content TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_payment_instructions_payment_method_id
        FOREIGN KEY (payment_method_id)
        REFERENCES payment_methods(id)
        ON DELETE CASCADE
);