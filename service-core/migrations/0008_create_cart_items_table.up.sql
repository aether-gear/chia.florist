CREATE TABLE cart_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    cart_id UUID NOT NULL,
    item_type VARCHAR(32) NOT NULL DEFAULT 'standard',
    product_id UUID,
    shop_id UUID NOT NULL,

    quantity INTEGER NOT NULL CHECK (quantity > 0),
    custom_design JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_cart
        FOREIGN KEY(cart_id)
        REFERENCES carts(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_product
        FOREIGN KEY(product_id)
        REFERENCES products(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_cart_item_shop
        FOREIGN KEY(shop_id)
        REFERENCES shops(id)
        ON DELETE RESTRICT,

    CONSTRAINT check_cart_items_type
        CHECK (
            (item_type = 'standard' AND product_id IS NOT NULL)
            OR (item_type = 'custom' AND product_id IS NULL AND custom_design IS NOT NULL)
        )
);

CREATE UNIQUE INDEX unique_standard_product_per_cart
ON cart_items(cart_id, product_id)
WHERE deleted_at IS NULL AND item_type = 'standard';

CREATE INDEX idx_cart_items_cart_id
ON cart_items(cart_id);

CREATE INDEX idx_cart_items_shop_id
ON cart_items(shop_id);
