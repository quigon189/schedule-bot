-- +goose Up
CREATE OR REPLACE VIEW group_schedules_view AS
SELECT
	sp.id as study_period_id,
	sp.academic_year,
	sp.half_year,  
	sp.start_date,
	sp.end_date,
	g.id as group_id,
	g.name as group_name,
	gs.semester,
	gs.schedule_image_url,
	gs.created_at
FROM group_schedules gs
JOIN study_periods sp ON gs.study_period_id = sp.id
JOIN groups g ON gs.group_id = g.id;

-- +goose Down
DROP VIEW group_schedules_view;
