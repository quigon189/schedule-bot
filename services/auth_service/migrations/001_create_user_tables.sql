-- +goose Up
CREATE TABLE users (
	telegram_id BIGINT PRIMARY KEY,
	username VARCHAR(255),
	full_name VARCHAR(255),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE roles (
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) UNIQUE NOT NULL,
	description TEXT
);

CREATE TABLE user_roles (
	user_id BIGINT REFERENCES users(telegram_id) ON DELETE CASCADE,
	role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
	PRIMARY KEY (user_id, role_id)
);

CREATE TABLE students (
	user_id BIGINT PRIMARY KEY REFERENCES users(telegram_id) ON DELETE CASCADE,
	group_name VARCHAR(100)
);

INSERT INTO roles (name, description) VALUES
	('admin', 'Administrator with full access'),
	('manager', 'Manager with limited administration access'),
	('teacher', 'Regular teacher user'),
	('student', 'Regular student user'),
	('user', 'New user');

-- +goose Down
DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;
