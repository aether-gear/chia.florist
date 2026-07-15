CREATE TABLE product_performance (
    product_id              UUID PRIMARY KEY,
    cost_price              BIGINT,
    supplier_lead_time_days INTEGER,
    gross_margin_pct        NUMERIC(6, 2),
    view_count              BIGINT NOT NULL DEFAULT 0,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ,

    CONSTRAINT fk_product_performance_product
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE CASCADE
);
