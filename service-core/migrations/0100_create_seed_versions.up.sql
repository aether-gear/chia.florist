CREATE TABLE IF NOT EXISTS seed_versions (
    name TEXT PRIMARY KEY,
    version TEXT,
    applied_at TIMESTAMPTZ DEFAULT NOW()
);