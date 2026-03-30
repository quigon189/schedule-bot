package repository

import (
	"context"
	"fmt"

	"core/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AcademicPeriodRepo struct {
	db *pgxpool.Pool
}

func NewAcademicPeriodRepo(db *pgxpool.Pool) *AcademicPeriodRepo {
	return &AcademicPeriodRepo{db: db}
}

// Create — создание учебного периода
func (r *AcademicPeriodRepo) Create(ctx context.Context, period *models.AcademicPeriod) error {
	query := `
	INSERT INTO schedule.academic_periods (year, semester, start_date, end_date)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		period.Year,
		period.Semester,
		period.StartDate,
		period.EndDate,
	).Scan(&period.ID, &period.CreatedAt, &period.UpdatedAt)
}

// GetByID — получение учебного периода по ID
func (r *AcademicPeriodRepo) GetByID(ctx context.Context, id int) (*models.AcademicPeriod, error) {
	var period models.AcademicPeriod
	query := `
	SELECT id, year, semester, start_date, end_date, created_at, updated_at
	FROM schedule.academic_periods
	WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
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
		return nil, fmt.Errorf("get academic period by id %d: %w", id, err)
	}
	return &period, nil
}

// GetActive — получение активного учебного периода
func (r *AcademicPeriodRepo) GetActive(ctx context.Context) (*models.AcademicPeriod, error) {
	var period models.AcademicPeriod
	query := `
	SELECT id, year, semester, start_date, end_date, created_at, updated_at
	FROM schedule.academic_periods
	WHERE CURRENT_DATE BETWEEN start_date AND end_date
	LIMIT 1
	`
	err := r.db.QueryRow(ctx, query).Scan(
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
		return nil, fmt.Errorf("get active academic period: %w", err)
	}
	return &period, nil
}

// GetAll — получение всех учебных периодов
func (r *AcademicPeriodRepo) GetAll(ctx context.Context) ([]models.AcademicPeriod, error) {
	query := `
	SELECT id, year, semester, start_date, end_date, created_at, updated_at
	FROM schedule.academic_periods
	ORDER BY year DESC, semester DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query academic periods: %w", err)
	}
	defer rows.Close()

	var periods []models.AcademicPeriod
	for rows.Next() {
		var period models.AcademicPeriod
		err := rows.Scan(
			&period.ID,
			&period.Year,
			&period.Semester,
			&period.StartDate,
			&period.EndDate,
			&period.CreatedAt,
			&period.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan academic period: %w", err)
		}
		periods = append(periods, period)
	}
	return periods, nil
}

// Update — обновление учебного периода
func (r *AcademicPeriodRepo) Update(ctx context.Context, period *models.AcademicPeriod) error {
	query := `
	UPDATE schedule.academic_periods
	SET year = $1, semester = $2, start_date = $3, end_date = $4, is_active = $5, updated_at = NOW()
	WHERE id = $6
	`
	_, err := r.db.Exec(ctx, query,
		period.Year,
		period.Semester,
		period.StartDate,
		period.EndDate,
		period.ID,
	)
	return err
}

// Delete — удаление учебного периода
func (r *AcademicPeriodRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM schedule.academic_periods WHERE id = $1`, id)
	return err
}

