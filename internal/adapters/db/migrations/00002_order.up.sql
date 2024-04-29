BEGIN TRANSACTION;

CREATE TYPE status_type AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE IF NOT EXISTS orders (
    id          bigint PRIMARY KEY,
    status      status_type NOT NULL DEFAULT 'NEW',
    user_id     bigint NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    accrual     double precision NULL,
    uploaded_at timestamp NOT NULL DEFAULT NOW()
);

COMMIT;
