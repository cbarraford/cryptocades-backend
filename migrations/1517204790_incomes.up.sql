CREATE TABLE incomes(
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    game_id INTEGER NOT NULL,
    session_id VARCHAR(12) NOT NULL,
    amount INTEGER NOT NULL,
    created_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX incomes_user_id_game_id_session_id ON incomes (user_id, game_id, session_id);

ALTER TABLE users DROP COLUMN mined_hashes, DROP COLUMN bonus_hashes;
