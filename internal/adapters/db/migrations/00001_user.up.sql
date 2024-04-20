BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

COMMIT;