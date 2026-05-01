CREATE TABLE shop_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    shop_id UUID NOT NULL,

    label TEXT NOT NULL,
    phone TEXT NOT NULL,
    
    is_active bool NOT NULL,

    province TEXT NOT NULL,
    city TEXT NOT NULL,
    district TEXT NOT NULL,
    village TEXT NOT NULL,
    full_address TEXT NOT NULL,
    postal_code TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_shop
        FOREIGN KEY(shop_id)
        REFERENCES shops(id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX uq_idx_one_active_address_per_shop
ON shop_addresses(shop_id)
WHERE is_active = true;