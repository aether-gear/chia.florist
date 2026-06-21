CREATE TABLE shop_couriers (
    shop_id UUID NOT NULL,

    code TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,

    PRIMARY KEY (shop_id, code),

    CONSTRAINT fk_shop_couriers_shop_id
        FOREIGN KEY (shop_id)
        REFERENCES shops(id)
        ON DELETE CASCADE,
);

CREATE INDEX idx_shop_couriers_shop_id
ON shop_couriers(shop_id);
