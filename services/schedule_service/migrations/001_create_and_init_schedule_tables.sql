-- +goose Up
CREATE TABLE groups (
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL UNIQUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESAMP
);

CREATE TABLE subjects (
	id SERIAL PRIMARY KEY,
	name VARCHAR(200) NOT NULL UNIQUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESAMP
);

CREATE TABLE classrooms (
	id SERIAL PRIMARY KEY,
	number VARCHAR(20) NOT NULL UNIQUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESAMP
);

CREATE TABLE study_periods (
	id SERIAL PRIMARY KEY,
	half_year INTEGER NOT NULL,
	academic_year VARCHAR(20) NOT NULL,
	start_date DATE NOT NULL,
	end_date DATE NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESAMP,

	UNIQUE(half_year, academic_year),
	CHECK(half_year BETWEEN 1 AND 2)
);

CREATE TABLE group_study_periods (
	study_period_id INTEGER REFERENCES study_periods(id),	
	group_id INTEGER REFERENCES groups(id),
	semester INTEGER NOT NULL UNIQUE,
	shedule_image_url VARCHAR(500),

	PRIMARY KEY(study_period_id, group_id)
);

CREATE TABLE schedule (
	id SERIAL PRIMARY KEY,
	study_period_id INTEGER NOT NULL REFERENCES study_periods(id) ON DELETE CASCADE,
	group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
	subject_id INTEGER NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
	classroom_id INTEGER NOT NULL REFERENCES classrooms(id) ON DELETE CASCADE,
	-- TODO: доделать!
);
