ALTER TABLE users ADD COLUMN confirmed BOOL NOT NULL DEFAULT FALSE;

CREATE TABLE confirmations(
    id SERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    email CITEXT NOT NULL,
    user_id int REFERENCES users(id) ON DELETE CASCADE,
    created_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);