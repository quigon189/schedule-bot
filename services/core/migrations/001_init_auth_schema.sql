-- +goose Up
CREATE SCHEMA auth;

CREATE TABLE auth.users (
	id SERIAL PRIMARY KEY,
	username VARCHAR(255) NOT NULL UNIQUE,
	full_name VARCHAR(255),
	email VARCHAR(255) NOT NULL UNIQUE,
	password_hash VARCHAR(255) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE auth.roles (
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) UNIQUE NOT NULL,
	description TEXT
);

CREATE TABLE auth.user_roles (
	user_id INTEGER REFERENCES auth.users(id) ON DELETE CASCADE,
	role_id INTEGER REFERENCES auth.roles(id) ON DELETE CASCADE,
	PRIMARY KEY (user_id, role_id)
);

CREATE TABLE auth.groups (
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL UNIQUE,
	specialty VARCHAR(255) NOT NULL,
	admission_year INTEGER NOT NULL
);

CREATE TABLE auth.student_profiels (
	user_id INTEGER PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
	group_id INTEGER REFERENCES auth.groups(id) ON DELETE SET NULL
);

CREATE TABLE auth.teacher_profiles (
	user_id INTEGER PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE
);

CREATE TABLE auth.subjects (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	group_id INTEGER NOT NULL REFERENCES auth.groups(id) ON DELETE CASCADE,
	teacher_id INTEGER REFERENCES auth.teacher_profiles(user_id) ON DELETE SET NULL,
	hours_load INTEGER NOT NULL DEFAULT 0,
	semester INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE auth.sessions (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id INTEGER NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
	refresh_token TEXT NOT NULL UNIQUE,
	user_agent TEXT,
	client_ip VARCHAR(45),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO auth.roles (name, description) VALUES
	('admin', 'Administrator with full access'),
	('manager', 'Manager with limited administration access'),
	('teacher', 'Regular teacher user'),
	('student', 'Regular student user'),
	('user', 'New user');

-- +goose Down
DROP TABLE IF EXISTS auth.sessions;
DROP TABLE IF EXISTS auth.subjects;
DROP TABLE IF EXISTS auth.teacher_profiles;
DROP TABLE IF EXISTS auth.student_profiels;
DROP TABLE IF EXISTS auth.groups;
DROP TABLE IF EXISTS auth.telegram_users;
DROP TABLE IF EXISTS auth.user_roles;
DROP TABLE IF EXISTS auth.roles;
DROP TABLE IF EXISTS auth.users;
