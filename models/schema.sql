CREATE TABLE IF NOT EXISTS guild (
    id bigint unsigned not null primary key,
    channel_id bigint unsigned not null,
    blocked boolean not null default false
);

CREATE INDEX blacklist_idx ON guild (blocked);

CREATE TABLE IF NOT EXISTS operator (
    id bigint unsigned not null primary key
);

