package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"core/internal/dto"
	"core/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubjectRepo struct {
	db                *pgxpool.Pool
	allowedSortFields map[string]bool
}

func NewSubjectRepo(db *pgxpool.Pool) *SubjectRepo {
	return &SubjectRepo{
		db: db,
		allowedSortFields: map[string]bool{
			"id":          true,
			"title":       true,
			"semester":    true,
			"hours_load":  true,
			"start_date":  true,
			"end_date":    true,
			"group_id":    true,
		},
	}
}

func (r *SubjectRepo) CreateSubject(ctx context.Context, subject *models.Subject) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := r.CreateSubjectWithTx(ctx, tx, subject); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *SubjectRepo) CreateSubjectWithTx(ctx context.Context, tx pgx.Tx, subject *models.Subject) error {
	query := `
	INSERT INTO schedule.subjects (title, semester, hours_load, start_date, end_date, group_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`
	err := tx.QueryRow(ctx, query,
		subject.Title,
		subject.Semester,
		subject.HoursLoad,
		subject.StartDate,
		subject.EndDate,
		subject.GroupID,
	).Scan(&subject.ID)
	return err
}

func (r *SubjectRepo) GetSubjectByID(ctx context.Context, id int) (*models.Subject, error) {
	subject := models.Subject{
		Group: models.Group{},
	}
	query := `
	SELECT s.id, s.title, s.semester, s.hours_load, s.start_date, s.end_date, s.group_id,
	       g.id, g.name, g.specialty, g.admission_year
	FROM schedule.subjects s
	LEFT JOIN auth.groups g ON s.group_id = g.id
	WHERE s.id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&subject.ID,
		&subject.Title,
		&subject.Semester,
		&subject.HoursLoad,
		&subject.StartDate,
		&subject.EndDate,
		&subject.GroupID,
		&subject.Group.ID,
		&subject.Group.Name,
		&subject.Group.Specialty,
		&subject.Group.AdmissionYear,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &subject, nil
}

func (r *SubjectRepo) GetAllSubjects(ctx context.Context, req *dto.PaginatedSubjectsRequest) (*dto.PaginatedSubjectsResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 10 // default
	}
	if perPage > 100 {
		perPage = 100
	}

	sortBy := req.SortBy
	if sortBy == "" || !r.allowedSortFields[sortBy] {
		sortBy = "id"
	}
	sortOrder := strings.ToUpper(req.SortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "ASC"
	}

	offset := (page - 1) * perPage
	orderClause := fmt.Sprintf("s.%s %s", sortBy, sortOrder)

	query := fmt.Sprintf(`
	SELECT s.id, s.title, s.semester, s.hours_load, s.start_date, s.end_date, s.group_id,
	       g.id, g.name, g.specialty, g.admission_year
	FROM schedule.subjects s
	LEFT JOIN auth.groups g ON s.group_id = g.id
	ORDER BY s.%s
	LIMIT $1 OFFSET $2
	`, orderClause)

	rows, err := r.db.Query(ctx, query, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("query subjects: %w", err)
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		subject := models.Subject{
			Group: models.Group{},
		}
		if err := rows.Scan(
			&subject.ID,
			&subject.Title,
			&subject.Semester,
			&subject.HoursLoad,
			&subject.StartDate,
			&subject.EndDate,
			&subject.GroupID,
			&subject.Group.ID,
			&subject.Group.Name,
			&subject.Group.Specialty,
			&subject.Group.AdmissionYear,
		); err != nil {
			return nil, fmt.Errorf("scan subject: %w", err)
		}
		subjects = append(subjects, subject)
	}

	var total int
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM schedule.subjects`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count subjects: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.PaginatedSubjectsResponse{
		Subjects:   subjects,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (r *SubjectRepo) GetSubjectsByDate(ctx context.Context, startDate, endDate time.Time) ([]models.Subject, error) {
	query := `
	SELECT s.id, s.title, s.semester, s.hours_load, s.start_date, s.end_date, s.group_id,
	       g.id, g.name, g.specialty, g.admission_year
	FROM schedule.subjects s
	LEFT JOIN auth.groups g ON s.group_id = g.id
	WHERE s.start_date >= $1 AND s.end_date <= $2
	ORDER BY g.name
	`
	rows, err := r.db.Query(ctx, query, startDate, endDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("subjects not found")
		}
		return nil, err
	}

	var subjects []models.Subject
	for rows.Next() {
		var subject models.Subject
		var group models.Group
		err := rows.Scan(
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
		)
		if err != nil {
			return nil, err
		}
		subject.Group = group
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func (r *SubjectRepo) GetSubjectsByGroupID(ctx context.Context, groupID int) ([]models.Subject, error) {
	query := `
	SELECT s.id, s.title, s.semester, s.hours_load, s.start_date, s.end_date, s.group_id,
	       g.id, g.name, g.specialty, g.admission_year
	FROM schedule.subjects s
	LEFT JOIN auth.groups g ON s.group_id = g.id
	WHERE s.group_id = $1
	ORDER BY s.semester
	`

	rows, err := r.db.Query(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("query subjects by group: %w", err)
	}
	defer rows.Close()

	var subjects []models.Subject
	for rows.Next() {
		subject := models.Subject{
			Group: models.Group{},
		}
		if err := rows.Scan(
			&subject.ID,
			&subject.Title,
			&subject.Semester,
			&subject.HoursLoad,
			&subject.StartDate,
			&subject.EndDate,
			&subject.GroupID,
			&subject.Group.ID,
			&subject.Group.Name,
			&subject.Group.Specialty,
			&subject.Group.AdmissionYear,
		); err != nil {
			return nil, fmt.Errorf("scan subject: %w", err)
		}
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

func (r *SubjectRepo) UpdateSubject(ctx context.Context, subject *models.Subject) error {
	query := `
	UPDATE schedule.subjects
	SET title = $1, semester = $2, hours_load = $3, start_date = $4, end_date = $5
	WHERE id = $6
	`
	_, err := r.db.Exec(ctx, query,
		subject.Title,
		subject.Semester,
		subject.HoursLoad,
		subject.StartDate,
		subject.EndDate,
		subject.ID,
	)
	return err
}

func (r *SubjectRepo) DeleteSubject(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM schedule.subjects WHERE id = $1`, id)
	return err
}
