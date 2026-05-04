CREATE TABLE cart_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    cart_id UUID NOT NULL,
    product_id UUID NOT NULL,
    shop_id UUID NOT NULL,

    quantity INTEGER NOT NULL CHECK (quantity > 0),

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
        ON DELETE RESTRICT
);

CREATE UNIQUE INDEX unique_product_per_carts
ON cart_items(cart_id, product_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_cart_items_cart_id
ON cart_items(cart_id);

CREATE INDEX idx_cart_items_shop_id
ON cart_items(shop_id);
