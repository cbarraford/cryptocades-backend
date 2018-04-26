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
    total_asteroids INTEGER NOT NULL DEFAULT 0,
    total_resources INTEGER NOT NULL DEFAULT 0,
    account_id INTEGER REFERENCES g2_accounts(id) ON DELETE CASCADE,
    session_id VARCHAR(12) DEFAULT '',
    health INTEGER NOT NULL DEFAULT 0,
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
    remaining INTEGER NOT NULL DEFAULT 0,
    distance INTEGER NOT NULL DEFAULT 0 CHECK (distance > 0),
    ship_id INTEGER NOT NULL DEFAULT 0,
    ship_speed INTEGER NOT NULL DEFAULT 1 CHECK (ship_speed > 0),
    solar_system INTEGER NOT NULL DEFAULT 0,
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
    created_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX g2_ship_upgrades_asset_id_category_id ON g2_ship_upgrades(asset_id, category_id);


-- 1 Engine
-- 2 Cargo
-- 3 Repair
-- 4 Hull
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (1, 1, 100, 100);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (1, 2, 200, 200);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (1, 3, 300, 300);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (1, 4, 400, 400);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (1, 5, 500, 500);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (1, 6, 600, 600);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (1, 7, 700, 700);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (1, 8, 800, 800);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (1, 9, 900, 900);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (1, 10, 1000, 1000);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (2, 1, 500, 100);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (2, 2, 700, 200);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (2, 3, 1000, 300);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (2, 4, 1200, 400);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (2, 5, 1500, 500);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (2, 6, 2000, 600);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (2, 7, 2300, 700);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (2, 8, 2500, 800);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (2, 9, 2700, 900);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (2, 10, 3000, 1000);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (3, 1, 1, 100);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (3, 2, 2, 200);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (3, 3, 3, 300);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (3, 4, 4, 400);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (3, 5, 5, 500);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (3, 6, 6, 600);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (3, 7, 7, 700);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (3, 8, 8, 800);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (3, 9, 9, 900);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (3, 10, 10, 1000);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (4, 1, 200, 100);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (4, 2, 400, 200);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (4, 3, 600, 300);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (4, 4, 800, 400);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (4, 5, 1000, 500);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (4, 6, 1200, 600);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (4, 7, 1400, 700);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (4, 8, 1600, 800);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (4, 9, 1800, 900);
INSERT INTO g2_ship_upgrades (category_id, asset_id, value, cost) VALUES 
    (4, 10, 2000, 1000);
