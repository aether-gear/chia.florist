CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    order_id UUID NOT NULL,
    product_variant_type product_variant_type NOT NULL DEFAULT 'standard',

    shop_id UUID NOT NULL,
    shop_name VARCHAR(255) NOT NULL,

    product_id UUID,
    product_name VARCHAR(255) NOT NULL,

    shipment_id UUID,
    courier_code VARCHAR(100),
    courier_service VARCHAR(100),
    shipping_fee_total BIGINT NOT NULL,

    quantity INTEGER NOT NULL,
    unit_price BIGINT NOT NULL,
    subtotal BIGINT NOT NULL,

    CONSTRAINT fk_order_items_order_id
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_order_items_shop_id
        FOREIGN KEY (shop_id)
        REFERENCES shops(id),

    CONSTRAINT fk_order_items_product_id
        FOREIGN KEY (product_id)
        REFERENCES products(id),

    CONSTRAINT check_order_items_type
        CHECK (
            (product_variant_type = 'standard' AND product_id IS NOT NULL)
            OR (product_variant_type = 'custom' AND product_id IS NULL)
        ),

    CONSTRAINT order_items_quantity_check
        CHECK (quantity > 0),

    CONSTRAINT order_items_unit_price_check
        CHECK (unit_price >= 0),

    CONSTRAINT order_items_subtotal_check
        CHECK (subtotal >= 0)
);

CREATE INDEX idx_order_items_shipment_id ON order_items(shipment_id);
