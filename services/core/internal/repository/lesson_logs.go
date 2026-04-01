package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"core/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LessonLogRepo struct {
	db *pgxpool.Pool
}

func NewLessonLogRepo(db *pgxpool.Pool) *LessonLogRepo {
	return &LessonLogRepo{db: db}
}

type LessonLogFilters struct {
	GroupID          *int       // идентификатор группы (через предмет)
	SubjectID        *int       // идентификатор предмета
	TeacherID        *int       // идентификатор преподавателя
	AudienceID       *int       // идентификатор аудитории
	AcademicPeriodID *int       // идентификатор учебного периода
	DateFrom         *time.Time // дата начала периода
	DateTo           *time.Time // дата окончания периода
	Status           *string    // статус занятия (planned, conducted, cancelled, etc.)
	Number           *int       // номер пары
}

func (r *LessonLogRepo) Create(ctx context.Context, log *models.LessonLog) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	err = r.CreateWithTx(ctx, tx, log)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *LessonLogRepo) CreateWithTx(ctx context.Context, tx pgx.Tx, log *models.LessonLog) error {
	query := `
	INSERT INTO schedule.lesson_logs (date, number, status, comment, subject_id, teacher_id, audience_id, academic_period_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		log.Date,
		log.Number,
		log.Status,
		log.Comment,
		log.SubjectID,
		log.TeacherID,
		log.AudienceID,
		log.AcademicPeriodID,
	).Scan(&log.ID)
}

// GetByID — получение записи журнала по ID с полной загрузкой связанных сущностей
func (r *LessonLogRepo) GetByID(ctx context.Context, id int) (*models.LessonLog, error) {
	query := `
	SELECT ll.id, ll.date, ll.number, ll.status, ll.comment, 
	       ll.subject_id, ll.teacher_id, ll.audience_id, ll.academic_period_id,
	       s.id, s.title, s.semester, s.hours_load, s.start_date, s.end_date, s.group_id,
	       g.id, g.name, g.specialty, g.admission_year,
	       t.id, t.username, t.full_name, t.email, t.created_at, t.updated_at,
	       a.id, a.name, a.number,
	       p.id, p.year, p.semester, p.start_date, p.end_date, p.created_at, p.updated_at
	FROM schedule.lesson_logs ll
	LEFT JOIN schedule.subjects s ON ll.subject_id = s.id
	LEFT JOIN auth.groups g ON s.group_id = g.id
	LEFT JOIN auth.users t ON ll.teacher_id = t.id
	LEFT JOIN schedule.audiences a ON ll.audience_id = a.id
	LEFT JOIN schedule.academic_periods p ON ll.academic_period_id = p.id
	WHERE ll.id = $1
	`

	var log models.LessonLog
	var subject models.Subject
	var group models.Group
	var teacher models.Teacher
	var user models.User
	var audience models.Audience
	var period models.AcademicPeriod

	err := r.db.QueryRow(ctx, query, id).Scan(
		&log.ID,
		&log.Date,
		&log.Number,
		&log.Status,
		&log.Comment,
		&log.SubjectID,
		&log.TeacherID,
		&log.AudienceID,
		&log.AcademicPeriodID,
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
		&period.ID,
		&period.Year,
		&period.Semester,
		&period.StartDate,
		&period.EndDate,
		&period.CreatedAt,
		&period.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get lesson log by id %d: %w", id, err)
	}

	// Заполняем связанные объекты
	subject.Group = group
	log.Subject = subject
	log.Teacher = teacher
	teacher.User = user
	log.Audience = audience
	log.AcademicPeriod = period

	return &log, nil
}

func (r *LessonLogRepo) GetAll(ctx context.Context, filters LessonLogFilters) ([]models.LessonLog, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	lessons, err := r.GetAllWithTx(ctx, tx, filters)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}
	return lessons, nil
}

// GetAll — получение списка записей журнала с фильтрацией
func (r *LessonLogRepo) GetAllWithTx(ctx context.Context, tx pgx.Tx, filters LessonLogFilters) ([]models.LessonLog, error) {
	var conditions []string
	var args []any
	argIndex := 1

	// Базовый запрос с JOIN
	query := `
	SELECT ll.id, ll.date, ll.number, ll.status, ll.comment, 
	       ll.subject_id, ll.teacher_id, ll.audience_id, ll.academic_period_id,
	       s.id, s.title, s.semester, s.hours_load, s.start_date, s.end_date, s.group_id,
	       g.id, g.name, g.specialty, g.admission_year,
	       t.id, t.username, t.full_name, t.email, t.created_at, t.updated_at,
	       a.id, a.name, a.number,
	       p.id, p.year, p.semester, p.start_date, p.end_date, p.created_at, p.updated_at
	FROM schedule.lesson_logs ll
	LEFT JOIN schedule.subjects s ON ll.subject_id = s.id
	LEFT JOIN auth.groups g ON s.group_id = g.id
	LEFT JOIN auth.users t ON ll.teacher_id = t.id
	LEFT JOIN schedule.audiences a ON ll.audience_id = a.id
	LEFT JOIN schedule.academic_periods p ON ll.academic_period_id = p.id
	WHERE 1=1
	`

	// Фильтр по группе (через group_id предмета)
	if filters.GroupID != nil {
		conditions = append(conditions, fmt.Sprintf("s.group_id = $%d", argIndex))
		args = append(args, *filters.GroupID)
		argIndex++
	}

	// Фильтр по предмету
	if filters.SubjectID != nil {
		conditions = append(conditions, fmt.Sprintf("ll.subject_id = $%d", argIndex))
		args = append(args, *filters.SubjectID)
		argIndex++
	}

	// Фильтр по преподавателю
	if filters.TeacherID != nil {
		conditions = append(conditions, fmt.Sprintf("ll.teacher_id = $%d", argIndex))
		args = append(args, *filters.TeacherID)
		argIndex++
	}

	// Фильтр по аудитории
	if filters.AudienceID != nil {
		conditions = append(conditions, fmt.Sprintf("ll.audience_id = $%d", argIndex))
		args = append(args, *filters.AudienceID)
		argIndex++
	}

	// Фильтр по учебному периоду
	if filters.AcademicPeriodID != nil {
		conditions = append(conditions, fmt.Sprintf("ll.academic_period_id = $%d", argIndex))
		args = append(args, *filters.AcademicPeriodID)
		argIndex++
	}

	// Фильтр по номеру пары
	if filters.Number != nil {
		conditions = append(conditions, fmt.Sprintf("ll.number = $%d", argIndex))
		args = append(args, *filters.Number)
		argIndex++
	}

	// Фильтр по статусу
	if filters.Status != nil {
		conditions = append(conditions, fmt.Sprintf("ll.status = $%d", argIndex))
		args = append(args, *filters.Status)
		argIndex++
	}

	// Фильтр по диапазону дат
	if filters.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("ll.date >= $%d", argIndex))
		args = append(args, *filters.DateFrom)
		argIndex++
	}
	if filters.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("ll.date <= $%d", argIndex))
		args = append(args, *filters.DateTo)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	// Сортировка по умолчанию: дата и номер пары
	query += " ORDER BY ll.date, ll.number"

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query lesson logs: %w", err)
	}
	defer rows.Close()

	var logs []models.LessonLog
	for rows.Next() {
		var log models.LessonLog
		var subject models.Subject
		var group models.Group
		var teacher models.Teacher
		var user models.User
		var audience models.Audience
		var period models.AcademicPeriod

		err := rows.Scan(
			&log.ID,
			&log.Date,
			&log.Number,
			&log.Status,
			&log.Comment,
			&log.SubjectID,
			&log.TeacherID,
			&log.AudienceID,
			&log.AcademicPeriodID,
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
			&period.ID,
			&period.Year,
			&period.Semester,
			&period.StartDate,
			&period.EndDate,
			&period.CreatedAt,
			&period.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan lesson log: %w", err)
		}

		subject.Group = group
		log.Subject = subject
		log.Teacher = teacher
		teacher.User = user
		log.Audience = audience
		log.AcademicPeriod = period

		logs = append(logs, log)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return logs, nil
}

// Update — обновление записи журнала
func (r *LessonLogRepo) Update(ctx context.Context, log *models.LessonLog) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	err = r.UpdateWithTx(ctx, tx, log)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *LessonLogRepo) UpdateWithTx(ctx context.Context, tx pgx.Tx, log *models.LessonLog) error {
	query := `
	UPDATE schedule.lesson_logs
	SET date = $1, number = $2, status = $3, comment = $4, 
	    subject_id = $5, teacher_id = $6, audience_id = $7, academic_period_id = $8
	WHERE id = $9
	`
	_, err := r.db.Exec(ctx, query,
		log.Date,
		log.Number,
		log.Status,
		log.Comment,
		log.SubjectID,
		log.TeacherID,
		log.AudienceID,
		log.AcademicPeriodID,
		log.ID,
	)
	return err
}

// Delete — удаление записи журнала по ID
func (r *LessonLogRepo) Delete(ctx context.Context, id int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	err = r.DeleteWithTx(ctx, tx, id)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *LessonLogRepo) DeleteWithTx(ctx context.Context, tx pgx.Tx, id int) error {
	_, err := tx.Exec(ctx, `DELETE FROM schedule.lesson_logs WHERE id = $1`, id)
	return err
}

// UpdateStatus — обновление статуса занятия
func (r *LessonLogRepo) UpdateStatus(ctx context.Context, id int, status, comment string) error {
	query := `
	UPDATE schedule.lesson_logs
	SET status = $1, comment = $2
	WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, status, comment, id)
	return err
}
