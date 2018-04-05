CREATE TABLE g2_accounts(
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    credits INTEGER NOT NULL DEFAULT 0 CHECK (credits >= 0),
    created_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE g2_ships(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL DEFAULT '',
    state INTEGER NOT NULL DEFAULT 0,
    total_asteroids INTEGER NOT NULL DEFAULT 0,
    total_resources INTEGER NOT NULL DEFAULT 0,
    account_id INTEGER REFERENCES g2_accounts(id) ON DELETE CASCADE,
    health INTEGER NOT NULL DEFAULT 0,
    drill_bit INTEGER NOT NULL DEFAULT 0,
    solar_system INTEGER NOT NULL DEFAULT 0,
    created_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE g2_applied_ship_upgrades(
    id SERIAL PRIMARY KEY,
    ship_id INTEGER REFERENCES g2_ships(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL DEFAULT 0,
    asset_id INTEGER NOT NULL DEFAULT 0,
    created_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX g2_ship_upgrades_ship_id_category_id ON g2_applied_ship_upgrades(ship_id, category_id);

CREATE TABLE g2_asteroids(
    id BIGSERIAL PRIMARY KEY,
    total INTEGER NOT NULL DEFAULT 0 CHECK (total > 0),
    remaining INTEGER NOT NULL DEFAULT 0 CHECK (remaining >= 0),
    distance INTEGER NOT NULL DEFAULT 0 CHECK (distance > 0),
    ship_id INTEGER NOT NULL DEFAULT 0,
    ship_speed INTEGER NOT NULL DEFAULT 1 CHECK (ship_speed > 0),
    solar_system INTEGER NOT NULL DEFAULT 0,
    session_id VARCHAR(12) DEFAULT '',
    created_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX g2_asteroids_ship_id ON g2_asteroids (ship_id)
WHERE ship_id > 0;

CREATE TABLE g2_logs(
    id BIGSERIAL PRIMARY KEY,
    ship_id INTEGER NOT NULL DEFAULT 0,
    level INTEGER NOT NULL DEFAULT 1,
    log TEXT NOT NULL DEFAULT '',
    created_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

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
    asteroid_id INTEGER REFERENCES g2_asteroids(id) ON DELETE CASCADE UNIQUE,
    session_id VARCHAR(12)
);

CREATE TABLE g2_ship_upgrades(
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL DEFAULT 0,
    asset_id INTEGER NOT NULL DEFAULT 0,
    cost INTEGER NOT NULL DEFAULT 0,
    value INTEGER NOT NULL DEFAULT 0,
    name TEXT,
    created_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX g2_ship_upgrades_asset_id_category_id ON g2_ship_upgrades(asset_id, category_id);


INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost, name) VALUES 
    (1, 1, 100, 100, 'Basic Engine');
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost, name) VALUES 
    (2, 1, 500, 100, '500 Cargo');
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost, name) VALUES 
    (3, 1, 100, 100, 'Copper Drill');
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost, name) VALUES 
    (4, 1, 100, 100, 'Copper Hull');
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost, name) VALUES 
    (4, 2, 200, 200, 'Aluminimu Hull');
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost, name) VALUES 
    (4, 3, 300, 300, 'Iron Hull');
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost, name) VALUES 
    (4, 4, 400, 400, 'Steel Hull');
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost, name) VALUES 
    (4, 5, 500, 500, 'Titanium Hull');
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost, name) VALUES 
    (4, 6, 600, 600, 'titanium Aluminide Hull');
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost, name) VALUES 
    (4, 7, 700, 700, 'Tungsten Hull');
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost, name) VALUES 
    (4, 8, 800, 800, 'Tungsten Carbide Hull');
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost, name) VALUES 
    (4, 9, 900, 900, 'Iconel Hull');
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost, name) VALUES 
    (4, 10, 1000, 1000, 'Carbon Hull');
