-- +goose Up
CREATE TABLE groups (
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL UNIQUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE subjects (
	id SERIAL PRIMARY KEY,
	name VARCHAR(200) NOT NULL UNIQUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE classrooms (
	id SERIAL PRIMARY KEY,
	number VARCHAR(20) NOT NULL UNIQUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE teachers (
	id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL UNIQUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE study_periods (
	id SERIAL PRIMARY KEY,
	half_year INTEGER NOT NULL,
	academic_year VARCHAR(20) NOT NULL,
	start_date DATE NOT NULL,
	end_date DATE NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

	UNIQUE(half_year, academic_year),
	CHECK(half_year BETWEEN 1 AND 2)
);

CREATE TABLE group_schedules (
	study_period_id INTEGER REFERENCES study_periods(id),	
	group_id INTEGER REFERENCES groups(id),
	semester INTEGER NOT NULL UNIQUE,
	schedule_image_url VARCHAR(500) UNIQUE,

	PRIMARY KEY(study_period_id, group_id),
	CHECK(semester BETWEEN 1 AND 10)
);

CREATE TABLE schedule (
	id SERIAL PRIMARY KEY,
	study_period_id INTEGER NOT NULL REFERENCES study_periods(id) ON DELETE CASCADE,
	group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
	subject_id INTEGER NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
	classroom_id INTEGER NOT NULL REFERENCES classrooms(id) ON DELETE CASCADE,
	teacher_id INTEGER NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
	
	day_of_week INTEGER NOT NULL ,
	lesson_number INTEGER NOT NULL,
	start_time TIME NOT NULL,
	end_time TIME NOT NULL,
	week_type INTEGER NOT NULL DEFAULT 0,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

	CHECK(day_of_week BETWEEN 1 AND 7),
	CHECK(lesson_number BETWEEN 0 AND 10),
	CHECK(week_type BETWEEN 0 AND 2),
	UNIQUE(study_period_id, day_of_week, lesson_number, group_id)
);

CREATE TABLE changes (
	id SERIAL PRIMARY KEY,
	change_date DATE NOT NULL,
	change_image_url VARCHAR(500) UNIQUE,
	change_description TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE schedule_changes (
	id SERIAL PRIMARY KEY,
	schedule_id INTEGER NOT NULL REFERENCES schedule(id) ON DELETE CASCADE,
	change_id INTEGER NOT NULL REFERENCES changes(id) ON DELETE CASCADE,
	
	substituted_subject_id INTEGER REFERENCES subjects(id) ON DELETE SET NULL,
	substituted_classroom_id INTEGER REFERENCES classrooms(id) ON DELETE SET NULL,
	substituted_teacher_id INTEGER REFERENCES teachers(id) ON DELETE SET NULL,

	is_cancelled BOOLEAN DEFAULT FALSE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_schedule_group_period ON schedule(group_id, study_period_id);
CREATE INDEX idx_schedule_day_lesson ON schedule(day_of_week, lesson_number);
CREATE INDEX idx_changes_date ON changes(change_date);

-- +goose Down
DROP TABLE IF EXISTS schedule_changes;
DROP TABLE IF EXISTS changes;
DROP TABLE IF EXISTS schedule;
DROP TABLE IF EXISTS group_schedules;
DROP TABLE IF EXISTS study_periods;
DROP TABLE IF EXISTS teachers;
DROP TABLE IF EXISTS classrooms;
DROP TABLE IF EXISTS subjects;
DROP TABLE IF EXISTS groups;
