BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    balance double precision NOT NULL DEFAULT 0,
    withdrawn double precision NOT NULL DEFAULT 0
);

COMMIT;
