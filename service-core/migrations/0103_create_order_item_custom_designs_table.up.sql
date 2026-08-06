CREATE TABLE order_item_custom_designs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_item_id UUID NOT NULL UNIQUE,
    version VARCHAR(16) NOT NULL DEFAULT '1.0.0',
    physical_size_id VARCHAR(64) NOT NULL,
    preview_url TEXT,

    header_text_upper VARCHAR(255),
    body_text_upper TEXT,
    header_text_lower VARCHAR(255),
    body_text_lower TEXT,

    design_snapshot JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_custom_design_order_item
        FOREIGN KEY (order_item_id)
        REFERENCES order_items(id)
        ON DELETE CASCADE
);
