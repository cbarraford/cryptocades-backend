ALTER TABLE users ADD COLUMN fb_id TEXT NOT NULL DEFAULT '';
ALTER TABLE users ALTER COLUMN username TYPE TEXT;
