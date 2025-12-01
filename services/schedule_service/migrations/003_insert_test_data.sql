-- +goose Up
INSERT INTO groups(name) VALUES
	('СА-501'),
	('СА-502'),
	('СА-503'),
	('Ш-206'),
	('Ш-208'),
	('Ш-210');

INSERT INTO study_periods(academic_year, half_year, start_date, end_date) VALUES
	('2025/2026', 1, '2025-09-01', '2025-12-27'),
	('2025/2026', 2, '2026-01-12', '2026-06-28');
