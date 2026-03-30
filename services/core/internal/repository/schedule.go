package repository

import (
	"context"
	"fmt"
	"strings"

	"core/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleRepo struct {
	db *pgxpool.Pool
}

func NewScheduleRepo(db *pgxpool.Pool) *ScheduleRepo {
	return &ScheduleRepo{db: db}
}

type ScheduleFilters struct {
	GroupID    *int // идентификатор группы (фильтр по group_id предмета)
	SubjectID  *int // идентификатор предмета
	TeacherID  *int // идентификатор преподавателя
	AudienceID *int // идентификатор аудитории
	DayOfWeek  *int // день недели (1–7)
	WeekType   *int // тип недели (0 – обе, 1 – числитель, 2 – знаменатель)
}

// Create — создание записи расписания
func (r *ScheduleRepo) Create(ctx context.Context, template *models.ScheduleTemplate) error {
	query := `
	INSERT INTO schedule.schedule_templates (day_of_week, number, week_type, subject_id, teacher_id, audience_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		template.DayOfWeek,
		template.Number,
		template.WeekType,
		template.SubjectID,
		template.TeacherID,
		template.AudienceID,
	).Scan(&template.ID)
}

// GetByID — получение записи расписания по ID с полной загрузкой связанных сущностей
func (r *ScheduleRepo) GetByID(ctx context.Context, id int) (*models.ScheduleTemplate, error) {
	query := `
	SELECT st.id, st.day_of_week, st.number, st.week_type, st.subject_id, st.teacher_id, st.audience_id,
	       s.id, s.title, s.semester, s.hours_load, s.start_date, s.end_date, s.group_id,
	       g.id, g.name, g.specialty, g.admission_year,
	       t.id, t.username, t.full_name, t.email, t.created_at, t.updated_at,
	       a.id, a.name, a.number
	FROM schedule.schedule_templates st
	LEFT JOIN schedule.subjects s ON st.subject_id = s.id
	LEFT JOIN auth.groups g ON s.group_id = g.id
	LEFT JOIN auth.users t ON st.teacher_id = t.id
	LEFT JOIN schedule.audiences a ON st.audience_id = a.id
	WHERE st.id = $1
	`

	var template models.ScheduleTemplate
	var subject models.Subject
	var group models.Group
	var teacher models.Teacher
	var user models.User
	var audience models.Audience

	err := r.db.QueryRow(ctx, query, id).Scan(
		&template.ID,
		&template.DayOfWeek,
		&template.Number,
		&template.WeekType,
		&template.SubjectID,
		&template.TeacherID,
		&template.AudienceID,
		&subject.ID,
		&subject.Title,
		&subject.Semester,
		&subject.HoursLoad,
		&subject.StartDate,
		&subject.EndDate,
		&subject.GroupID,
		&group.ID,
		&group.Name,
		&group.Specialty,
		&group.AdmissionYear,
		&user.ID,
		&user.Name,
		&user.FullName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&audience.ID,
		&audience.Name,
		&audience.Number,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get schedule by id %d: %w", id, err)
	}

	subject.Group = group
	template.Subject = subject

	teacher.User = user
	template.Teacher = teacher

	template.Audience = audience

	return &template, nil
}

// GetAll — получение списка записей расписания с фильтрацией и привязкой к семестру
// Параметр semester позволяет ограничить выборку предметами указанного семестра.
// Если semester == nil, фильтр по семестру не применяется.
func (r *ScheduleRepo) GetAll(ctx context.Context, filters ScheduleFilters, semester *int) ([]models.ScheduleTemplate, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Базовый запрос с JOIN
	query := `
	SELECT st.id, st.day_of_week, st.number, st.week_type, st.subject_id, st.teacher_id, st.audience_id,
	       s.id, s.title, s.semester, s.hours_load, s.start_date, s.end_date, s.group_id,
	       g.id, g.name, g.specialty, g.admission_year,
	       t.id, t.username, t.full_name, t.email, t.created_at, t.updated_at,
	       a.id, a.name, a.number
	FROM schedule.schedule_templates st
	LEFT JOIN schedule.subjects s ON st.subject_id = s.id
	LEFT JOIN auth.groups g ON s.group_id = g.id
	LEFT JOIN auth.users t ON st.teacher_id = t.id
	LEFT JOIN schedule.audiences a ON st.audience_id = a.id
	`

	// Фильтр по семестру
	if semester != nil {
		conditions = append(conditions, fmt.Sprintf("s.semester = $%d", argIndex))
		args = append(args, *semester)
		argIndex++
	}

	// Фильтр по группе (через group_id предмета)
	if filters.GroupID != nil {
		conditions = append(conditions, fmt.Sprintf("s.group_id = $%d", argIndex))
		args = append(args, *filters.GroupID)
		argIndex++
	}

	// Фильтр по предмету
	if filters.SubjectID != nil {
		conditions = append(conditions, fmt.Sprintf("st.subject_id = $%d", argIndex))
		args = append(args, *filters.SubjectID)
		argIndex++
	}

	// Фильтр по преподавателю
	if filters.TeacherID != nil {
		conditions = append(conditions, fmt.Sprintf("st.teacher_id = $%d", argIndex))
		args = append(args, *filters.TeacherID)
		argIndex++
	}

	// Фильтр по аудитории
	if filters.AudienceID != nil {
		conditions = append(conditions, fmt.Sprintf("st.audience_id = $%d", argIndex))
		args = append(args, *filters.AudienceID)
		argIndex++
	}

	// Фильтр по дню недели
	if filters.DayOfWeek != nil {
		conditions = append(conditions, fmt.Sprintf("st.day_of_week = $%d", argIndex))
		args = append(args, *filters.DayOfWeek)
		argIndex++
	}

	// Фильтр по типу недели
	if filters.WeekType != nil {
		conditions = append(conditions, fmt.Sprintf("st.week_type = $%d", argIndex))
		args = append(args, *filters.WeekType)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	// Сортировка по умолчанию: день недели и номер пары
	query += " ORDER BY st.day_of_week, st.number"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query schedule templates: %w", err)
	}
	defer rows.Close()

	var templates []models.ScheduleTemplate
	for rows.Next() {
		var template models.ScheduleTemplate
		var subject models.Subject
		var group models.Group
		var teacher models.Teacher
		var user models.User
		var audience models.Audience

		err := rows.Scan(
			&template.ID,
			&template.DayOfWeek,
			&template.Number,
			&template.WeekType,
			&template.SubjectID,
			&template.TeacherID,
			&template.AudienceID,
			&subject.ID,
			&subject.Title,
			&subject.Semester,
			&subject.HoursLoad,
			&subject.StartDate,
			&subject.EndDate,
			&subject.GroupID,
			&group.ID,
			&group.Name,
			&group.Specialty,
			&group.AdmissionYear,
			&user.ID,
			&user.Name,
			&user.FullName,
			&user.Email,
			&user.PasswordHash,
			&user.CreatedAt,
			&user.UpdatedAt,
			&audience.ID,
			&audience.Name,
			&audience.Number,
		)
		if err != nil {
			return nil, fmt.Errorf("scan schedule template: %w", err)
		}

		subject.Group = group
		template.Subject = subject

		teacher.User = user
		template.Teacher = teacher

		template.Audience = audience

		templates = append(templates, template)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return templates, nil
}

// Update — обновление существующей записи расписания
func (r *ScheduleRepo) Update(ctx context.Context, template *models.ScheduleTemplate) error {
	query := `
	UPDATE schedule.schedule_templates
	SET day_of_week = $1, number = $2, week_type = $3, subject_id = $4, teacher_id = $5, audience_id = $6
	WHERE id = $7
	`
	_, err := r.db.Exec(ctx, query,
		template.DayOfWeek,
		template.Number,
		template.WeekType,
		template.SubjectID,
		template.TeacherID,
		template.AudienceID,
		template.ID,
	)
	return err
}

// Delete — удаление записи расписания по ID
func (r *ScheduleRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM schedule.schedule_templates WHERE id = $1`, id)
	return err
}
