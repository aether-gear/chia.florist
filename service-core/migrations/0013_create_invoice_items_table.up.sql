CREATE TABLE invoice_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    invoice_id UUID NOT NULL,
    product_id UUID NOT NULL,

    product_name TEXT NOT NULL,
    product_price NUMERIC(15, 2) NOT NULL,

    quantity INTEGER NOT NULL,
    subtotal NUMERIC(15, 2) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_invoice
        FOREIGN KEY(invoice_id)
        REFERENCES invoices(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_product
        FOREIGN KEY(product_id)
        REFERENCES products(id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_product_per_invoice
ON invoice_items(invoice_id, product_id);