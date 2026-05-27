package repository

import (
	"context"
	"core/internal/models"
	"errors"
	"fmt"

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
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := r.CreateWithTx(ctx, tx, group); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
func (r *GroupRepo) CreateWithTx(ctx context.Context, tx pgx.Tx, group *models.Group) error {
	query := `
	INSERT INTO auth.groups (name, specialty, admission_year)
	VALUES ($1, $2, $3)
	RETURNING id
	`
	return tx.QueryRow(ctx, query, group.Name, group.Specialty, group.AdmissionYear).Scan(&group.ID)
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

	query = `
	SELECT u.id, u.username, u.full_name, u.email
	FROM auth.users u
	JOIN auth.student_profiles sp ON u.id = sp.user_id
	WHERE sp.group_id = $1
	ORDER BY u.full_name
	`
	rows, err := r.db.Query(ctx, query, group.ID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return &group, nil
	}

	for rows.Next() {
		var student models.User
		if err := rows.Scan(
			&student.ID,
			&student.Name,
			&student.FullName,
			&student.Email,
		); err != nil {
			return &group, nil
		}
		group.Students = append(group.Students, student)
	}

	return &group, nil
}

func (r *GroupRepo) GetByName(ctx context.Context, name string) (*models.Group, error) {
	var group models.Group
	query := `
	SELECT id, name, specialty, admission_year
	FROM auth.groups
	WHERE name LIKE '%' || $1 || '%'
	`
	err := r.db.QueryRow(ctx, query, name).Scan(
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

	query = `
	SELECT u.id, u.username, u.full_name, u.email
	FROM auth.users u
	JOIN auth.student_profiles sp ON u.id = sp.user_id
	WHERE sp.group_id = $1
	ORDER BY u.full_name
	`
	for i := range groups {
		rows, err := r.db.Query(ctx, query, groups[i].ID)
		if err != nil {
			continue
		}

		for rows.Next() {
			var student models.User
			if err := rows.Scan(
				&student.ID,
				&student.Name,
				&student.FullName,
				&student.Email,
			); err != nil {
				continue
			}

			groups[i].Students = append(groups[i].Students, student)
		}

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
