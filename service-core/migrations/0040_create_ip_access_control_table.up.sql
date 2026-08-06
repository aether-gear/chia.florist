CREATE TYPE ip_status AS ENUM ('banned', 'whitelisted', 'ignored', 'banned_muted', 'whitelisted_muted');

CREATE TABLE ip_access_control (
    ip          TEXT        PRIMARY KEY,
    status      ip_status   NOT NULL,
    reason      TEXT        NOT NULL DEFAULT '',
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ip_access_control_status
    ON ip_access_control(status);
