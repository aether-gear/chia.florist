CREATE TYPE audit_category
    AS ENUM (
        'user_action',
        'waf_event',
        'threat_intel',
        'system_job'
    );

CREATE TYPE audit_outcome
    AS ENUM (
        'success',
        'failure',
        'blocked'
    );

CREATE TABLE audit_logs (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),

    category    audit_category  NOT NULL,
    action      TEXT            NOT NULL,
    resource    TEXT            NOT NULL,
    resource_id TEXT,

    actor_id    TEXT,
    outcome     audit_outcome   NOT NULL,

    request_id  TEXT,
    client_ip   TEXT,

    metadata    JSONB,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_category
    ON audit_logs(category);

CREATE INDEX idx_audit_logs_actor_id
    ON audit_logs(actor_id);

CREATE INDEX idx_audit_logs_request_id
    ON audit_logs(request_id);

CREATE INDEX idx_audit_logs_action
    ON audit_logs(action);

CREATE INDEX idx_audit_logs_created_at
    ON audit_logs(created_at DESC);
