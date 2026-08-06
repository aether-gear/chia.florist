CREATE TABLE invoice_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL,
    product_variant_type product_variant_type NOT NULL DEFAULT 'standard',

    shop_id UUID NOT NULL,
    shop_name VARCHAR(255) NOT NULL,

    product_id UUID,
    product_name VARCHAR(255) NOT NULL,

    quantity INTEGER NOT NULL,
    unit_price BIGINT NOT NULL,
    subtotal BIGINT NOT NULL,

    courier_code VARCHAR(100),
    courier_service VARCHAR(100),
    shipping_fee_total BIGINT NOT NULL DEFAULT 0,

    CONSTRAINT fk_invoice_items_invoice_id
        FOREIGN KEY (invoice_id)
        REFERENCES invoices(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_invoice_items_shop_id
        FOREIGN KEY (shop_id)
        REFERENCES shops(id),

    CONSTRAINT fk_invoice_items_product_id
        FOREIGN KEY (product_id)
        REFERENCES products(id),

    CONSTRAINT invoice_items_quantity_check
        CHECK (quantity > 0),

    CONSTRAINT invoice_items_unit_price_check
        CHECK (unit_price >= 0),

    CONSTRAINT invoice_items_subtotal_check
        CHECK (subtotal >= 0)
);
