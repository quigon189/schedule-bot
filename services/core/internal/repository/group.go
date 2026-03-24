package repository

import (
	"context"
	"core/internal/models"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GroupRepo struct {
	db *pgxpool.Pool
}

func NewGroupRepo(db *pgxpool.Pool) *GroupRepo {
	return &GroupRepo{db: db}
}

func (r *GroupRepo) Create(ctx context.Context, group *models.Group) error {
	query := `
	INSERT INTO auth.groups (name, specialty, admission_year)
	VALUES ($1, $2, $3)
	RETURNING id
	`
	return r.db.QueryRow(ctx, query, group.Name, group.Specialty, group.AdmissionYear).Scan(&group.ID)
}

func (r *GroupRepo) GetByID(ctx context.Context, id int) (*models.Group, error) {
	var group models.Group
	query := `
	SELECT id, name, specialty, admission_year
	FROM auth.groups
	WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&group.ID,
		&group.Name,
		&group.Specialty,
		&group.AdmissionYear,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}

func (r *GroupRepo) GetAll(ctx context.Context) ([]models.Group, error) {
	groups := []models.Group{}
	query := `
	SELECT id, name, specialty, admission_year
	FROM auth.groups
	ORDER BY id
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var g models.Group
		if err := rows.Scan(&g.ID, &g.Name, &g.Specialty, &g.AdmissionYear); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func (r *GroupRepo) Update(ctx context.Context, group *models.Group) error {
	query := `
	UPDATE auth.groups
	SET name = $1, specialty = $2, admission_year = $3
	WHERE id = $4
	`
	_, err := r.db.Exec(ctx, query, group.Name, group.Specialty, group.AdmissionYear, group.ID)
	return err
}

func (r *GroupRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM auth.groups WHERE id = $1`, id)
	return err
}
