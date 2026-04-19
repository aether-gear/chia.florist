CREATE TABLE addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL,

    recipient_name TEXT NOT NULL,
    phone TEXT NOT NULL,
    
    is_default bool NOT NULL,

    province TEXT NOT NULL,
    city TEXT NOT NULL,
    district TEXT NOT NULL,
    village TEXT NOT NULL,
    full_address TEXT NOT NULL,
    postal_code TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX uq_idx_one_default_address_per_user
ON addresses(user_id)
WHERE is_default = true;