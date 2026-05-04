CREATE TABLE transaction_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    transaction_id UUID NOT NULL,
    product_id UUID NOT NULL,
    shop_id UUID NOT NULL,

    quantity INTEGER NOT NULL CHECK (quantity > 0),

    product_name TEXT NOT NULL,
    product_price NUMERIC(15, 2) NOT NULL,

    subtotal NUMERIC(15, 2) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_transaction
        FOREIGN KEY(transaction_id)
        REFERENCES transactions(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_product
        FOREIGN KEY(product_id)
        REFERENCES products(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_transaction_item_shop
        FOREIGN KEY(shop_id)
        REFERENCES shops(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_transaction_item_single_shop
        FOREIGN KEY(transaction_id, shop_id)
        REFERENCES transactions(id, shop_id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_product_shop_per_transaction
ON transaction_items(transaction_id, product_id, shop_id);