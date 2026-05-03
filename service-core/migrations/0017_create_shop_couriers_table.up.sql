CREATE TABLE shop_couriers (
    shop_id UUID NOT NULL,
    
    code TEXT NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,

    PRIMARY KEY (shop_id, code)
);

CREATE INDEX idx_shop_couriers_shop_id
ON shop_couriers(shop_id);