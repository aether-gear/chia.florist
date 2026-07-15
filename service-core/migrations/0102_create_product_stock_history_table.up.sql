CREATE TABLE product_stock_history (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id  UUID        NOT NULL,
    shop_id     UUID        NOT NULL,
    available   INTEGER     NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_psh_product
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_psh_shop
        FOREIGN KEY (shop_id)
        REFERENCES shops(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_product_stock_history_product_id
    ON product_stock_history(product_id);

CREATE INDEX idx_product_stock_history_recorded_at
    ON product_stock_history(recorded_at DESC);
