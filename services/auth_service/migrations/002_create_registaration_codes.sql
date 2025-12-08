-- +goose Up
CREATE TABLE registration_codes (
	id SERIAL PRIMARY KEY,
	code VARCHAR(20) UNIQUE NOT NULL,
	role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
	group_name VARCHAR(20),
	max_uses INT NOT NULL DEFAULT 1,
	current_uses INT NOT NULL DEFAULT 0,
	created_by INT REFERENCES users(telegram_id) ON DELETE SET NULL,
	expires_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_registration_codes_code ON registration_codes(code);
CREATE INDEX idx_registration_codes_created_by ON registration_codes(created_by);
CREATE INDEX idx_registration_codes_expires_at ON registration_codes(expires_at);

CREATE OR REPLACE FUNCTION delete_expired_codes()
RETURNS void AS $$
BEGIN
	DELETE FROM registration_codes
	WHERE expires_at < NOW()
	   OR current_uses >= max_uses;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_check_expiration
	BEFORE INSERT OR UPDATE registration_codes
	FOR EACH ROW
	EXECUTE delete_expired_codes();

-- +goose Down
DROP TABLE IF EXISTS registration_codes;
DROP TRIGGER IF EXISTS trigger_check_expiration ON registration_codes;
DROP FUNCTION IF EXISTS delete_expired_codes();
