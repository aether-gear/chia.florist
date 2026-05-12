CREATE TABLE product_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,

    thumbnail_url TEXT NOT NULL,
    preview_url   TEXT NOT NULL,
    detail_url    TEXT NOT NULL,

    thumbnail_key TEXT NOT NULL,
    preview_key   TEXT NOT NULL,
    detail_key    TEXT NOT NULL,

    is_primary BOOLEAN NOT NULL DEFAULT false,
    display_order INT NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    UNIQUE(product_id, display_order)
);