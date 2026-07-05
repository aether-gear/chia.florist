CREATE TABLE waf_rules (
    id          TEXT        PRIMARY KEY,
    description TEXT        NOT NULL,
    pattern     TEXT        NOT NULL,
    tags        TEXT[]      NOT NULL DEFAULT '{}',
    impact      TEXT        NOT NULL DEFAULT '',
    enabled     BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_waf_rules_enabled
    ON waf_rules(enabled);
