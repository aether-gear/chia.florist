CREATE TYPE shop_approval_status AS ENUM ('pending', 'approved', 'rejected');

CREATE TABLE shops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT,

    is_active bool NOT NULL,
    approval_status shop_approval_status NOT NULL DEFAULT 'pending',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);
