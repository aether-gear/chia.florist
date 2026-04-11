CREATE TYPE product_status AS ENUM ('active', 'inactive', 'archived');

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    sku TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT,
    status product_status DEFAULT 'active',

    base_price NUMERIC(15, 2) NOT NULL,
    weight NUMERIC(10, 2),

    stock INTEGER NOT NULL DEFAULT 0,
    reserved_stock INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT check_stock_non_negative
        CHECK (stock >= 0 AND reserved_stock >= 0)
);