-- +goose Up
CREATE OR REPLACE VIEW schedule_with_changes AS
SELECT
	sp.academic_year,
	sp.half_year,
	gs.semester,
	g.name as group_name,
	s.day_of_week,
	s.lesson_number,
	s.start_time,
	s.end_time,
	
	sub_orig.name as original_subject,
	c_orig.number as original_classroom,
	t_orig.name as original_teacher,

	ch.change_date,
	sub_sub.name as substituted_subject,
	c_sub.number as substituted_classroom,
	t_sub.name as substituted_teacher,
	sc.is_cancelled,
	ch.change_image_url,
	ch.change_description,
	gs.schedule_image_url as group_schedule_image_url
FROM schedule s
JOIN study_periods sp ON s.study_period_id = sp.id
JOIN group_schedules gs ON s.study_period_id = gs.study_period_id AND s.group_id = gs.group_id
JOIN groups g ON s.group_id = g.id
JOIN subjects sub_orig ON s.subject_id = sub_orig.id
JOIN classrooms c_orig ON s.classroom_id = c_orig.id
JOIN teachers t_orig ON s.teacher_id = t_orig.id
LEFT JOIN schedule_changes sc ON s.id = sc.schedule_id
LEFT JOIN changes ch ON sc.change_id = ch.id
LEFT JOIN subjects sub_sub ON sc.substituted_subject_id = sub_sub.id
LEFT JOIN classrooms c_sub ON sc.substituted_classroom_id = c_sub.id
LEFT JOIN teachers t_sub ON sc.substituted_teacher_id = t_sub.id
ORDER BY sp.academic_year DESC, sp.half_year DESC, g.name, s.day_of_week, s.lesson_number;
