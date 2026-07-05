CREATE TYPE filter_entry_type AS ENUM ('keyword', 'url');

CREATE TABLE filter_config (
    value   TEXT                NOT NULL,
    type    filter_entry_type   NOT NULL,
    PRIMARY KEY (value, type)
);
