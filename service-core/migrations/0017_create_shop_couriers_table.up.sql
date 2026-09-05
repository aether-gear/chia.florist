CREATE TABLE shop_couriers (
    shop_id UUID NOT NULL,

    code TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT FALSE,
    branch_name VARCHAR(150),
    location_address TEXT,

    verification_status VARCHAR(30) NOT NULL DEFAULT 'unconfigured',
    verified_at TIMESTAMPTZ,
    verified_by UUID REFERENCES staff(id) ON DELETE SET NULL,
    rejection_reason TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    PRIMARY KEY (shop_id, code),

    CONSTRAINT fk_shop_couriers_shop_id
        FOREIGN KEY (shop_id)
        REFERENCES shops(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_shop_couriers_active_details
        CHECK (
            active = FALSE OR (
                branch_name IS NOT NULL AND BTRIM(branch_name) <> '' AND
                location_address IS NOT NULL AND BTRIM(location_address) <> ''
            )
        )
);

CREATE INDEX idx_shop_couriers_shop_id
ON shop_couriers(shop_id);

CREATE INDEX idx_shop_couriers_shop_verification
ON shop_couriers(shop_id, verification_status);

CREATE INDEX idx_shop_couriers_active_details
ON shop_couriers(shop_id, code)
WHERE active = TRUE;
