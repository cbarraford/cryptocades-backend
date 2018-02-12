CREATE FUNCTION make_code() RETURNS text AS $$
DECLARE
    code text;
    done bool;
BEGIN
    done := false;
    WHILE NOT done LOOP
        code := md5(''||now()::text||random()::text);
        done := NOT exists(SELECT 1 FROM users WHERE referral_code=code);
    END LOOP;
    RETURN code;
END;
$$ LANGUAGE PLPGSQL VOLATILE;

ALTER TABLE users ADD COLUMN referral_code TEXT UNIQUE NOT NULL DEFAULT make_code();

ALTER TABLE incomes ALTER COLUMN session_id TYPE TEXT;
