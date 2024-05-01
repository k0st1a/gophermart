BEGIN;

CREATE TABLE IF NOT EXISTS withdrawals (
    id          bigserial PRIMARY KEY,
    order_id    bigint NOT NULL,
    user_id     bigint NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    sum         double precision NOT NULL,
    processed_at timestamp NOT NULL DEFAULT NOW()
);

COMMIT;
