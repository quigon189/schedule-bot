-- +goose UP
CREATE SCHEMA schedule;

CREATE TYPE lesson_status AS ENUM ('planned', 'completed', 'canceled', 'rescheduled');

CREATE TABLE schedule.academic_periods (
	id SERIAL PRIMARY KEY,
	year VARCHAR(10) NOT NULL,
	semester INTEGER NOT NULL CHECK(semester IN (1,2)),
	start_date DATE NOT NULL,
	end_date DATE NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(year, semester)
);

CREATE TABLE schedule.subjects (
	id SERIAL PRIMARY KEY,
	title VARCHAR(255) NOT NULL,
	semester INTEGER NOT NULL,
	group_id INTEGER REFERENCES auth.groups(id) ON DELETE CASCADE,
	hours_load INTEGER NOT NULL DEFAULT 0,
	start_date DATE NOT NULL,
	end_date DATE NOT NULL,

	CHECK(semester BETWEEN 1 AND 10)
);

CREATE TABLE schedule.audiences (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	number VARCHAR(255) NOT NULL
);

CREATE TABLE schedule.templates (
	id SERIAL PRIMARY KEY,
	day_of_week INTEGER,
	number INTEGER ,
	week_type INTEGER,
	academic_period_id INTEGER REFERENCES schedule.academic_periods(id) ON DELETE CASCADE,
	subject_id INTEGER REFERENCES schedule.subjects(id) ON DELETE CASCADE,
	teacher_id INTEGER REFERENCES auth.teacher_profiles(user_id) ON DELETE SET NULL,
	audience_id INTEGER REFERENCES schedule.audiences(id) ON DELETE SET NULL,

	CHECK(week_type BETWEEN 0 AND 2),
    CHECK(number BETWEEN 1 AND 7),
    CHECK(day_of_week BETWEEN 1 AND 7)
);

CREATE TABLE schedule.lesson_logs (
	id SERIAL PRIMARY KEY,
	date DATE NOT NULL,
	number INTEGER,
	academic_period_id INTEGER REFERENCES schedule.academic_periods(id) ON DELETE CASCADE,
	subject_id INTEGER REFERENCES schedule.subjects(id) ON DELETE CASCADE,
	teacher_id INTEGER REFERENCES auth.teacher_profiles(user_id) ON DELETE SET NULL,
	audience_id INTEGER REFERENCES schedule.audiences(id) ON DELETE SET NULL,
	status lesson_status DEFAULT 'planned',
	comment TEXT
);
