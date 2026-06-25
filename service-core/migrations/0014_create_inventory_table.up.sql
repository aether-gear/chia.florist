CREATE TABLE inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    product_id UUID NOT NULL,
    shop_id UUID NOT NULL,

    stock INTEGER NOT NULL DEFAULT 0,
    reserved_stock INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT fk_product
        FOREIGN KEY(product_id)
        REFERENCES products(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_inventory_shop
        FOREIGN KEY(shop_id)
        REFERENCES shops(id)
        ON DELETE CASCADE,

    CONSTRAINT unique_inventory_product_shop
        UNIQUE(product_id, shop_id)
);

CREATE INDEX idx_inventory_shop_id
ON inventory(shop_id);