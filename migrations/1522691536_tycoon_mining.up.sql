CREATE TABLE g2_ledgers(
    id BIGSERIAL PRIMARY KEY,
    account_id INTEGER REFERENCES g2_accounts(id) ON DELETE CASCADE, 
    session_id VARCHAR(12),
    amount INTEGER NOT NULL,
    description TEXT,
    created_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX g2_mining_account_id_session_id ON g2_ledgers (account_id, session_id);

CREATE TABLE g2_sessions(
    id BIGSERIAL PRIMARY KEY,
    ship_id INTEGER REFERENCES g2_ships(id) ON DELETE CASCADE UNIQUE,
    session_id VARCHAR(12),
    touched TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
