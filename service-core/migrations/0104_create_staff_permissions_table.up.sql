CREATE TABLE staff_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id UUID NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
    shop_id UUID NOT NULL REFERENCES shops(id) ON DELETE CASCADE,

    permissions JSONB NOT NULL DEFAULT '[]'::jsonb,
    rules JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT uq_staff_shop
        UNIQUE (staff_id, shop_id)
);

CREATE INDEX idx_staff_permissions_staff_id
    ON staff_permissions(staff_id);

CREATE INDEX idx_staff_permissions_shop_id
    ON staff_permissions(shop_id);
