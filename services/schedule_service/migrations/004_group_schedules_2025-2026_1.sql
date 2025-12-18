-- +goose Up
INSERT INTO groups(name) VALUES
	('СА-501'),
	('СА-502'),
	('СА-503'),
	('Ш-206')
ON CONFLICT (name)
DO NOTHING;

INSERT INTO study_periods(academic_year, half_year, start_date, end_date) VALUES
	('2025/2026', 1, '2025-09-01', '2025.12.28')
ON CONFLICT (academic_year, half_year)
DO NOTHING;
	
INSERT INTO group_schedules(study_period_id, group_id, semester, schedule_image_url) VALUES
	(
		(SELECT id FROM study_periods WHERE academic_year = '2025/2026' AND half_year = 1),
		(SELECT id FROM groups WHERE name = 'СА-501'),
		5,
		'https://r1.nubex.ru/s824-dc0/f14600_bf/%D0%A1%D0%90-501.jpg'
	),
	(
		(SELECT id FROM study_periods WHERE academic_year = '2025/2026' AND half_year = 1),
		(SELECT id FROM groups WHERE name = 'СА-502'),
		5,
		'https://r1.nubex.ru/s824-dc0/f14601_44/%D0%A1%D0%90-502.jpg'
	)
ON CONFLICT (study_period_id, group_id)
DO NOTHING;
