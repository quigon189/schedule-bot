package repository

import (
	"context"
	"core/internal/models"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepo struct {
	db *pgxpool.Pool
}

func NewRoleRepo(db *pgxpool.Pool) *RoleRepo {
	return &RoleRepo{db: db}
}

func (r *RoleRepo) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	query := `SELECT id, name, description FROM auth.roles WHERE name = $1`
	err := r.db.QueryRow(ctx, query, name).Scan(&role.ID, &role.Name, &role.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) AssignRoleToUser(ctx context.Context, userID int, roleName string) error {
	query := `
	INSERT INTO auth.user_roles (user_id, role_id)
	VALUES ($1, (SELECT id FROM auth.roles WHERE name = $2))
	ON CONFLICT (user_id, role_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, userID, roleName)
	return err
}

func (r *RoleRepo) RemoveRoleFromUser(ctx context.Context, userID int, roleName string) error {
	query := `
	DELETE FROM auth.user_roles
	WHERE user_id = $1 AND role_id = (SELECT id FROM auth.roles WHERE name = $2)
	`
	_, err := r.db.Exec(ctx, query, userID, roleName)
	return err
}

func (r *RoleRepo) GetUserRoles(ctx context.Context, userID int) ([]models.Role, error) {
	query := `
	SELECT r.id, r.name, r.description
	FROM auth.roles r
	JOIN auth.user_roles ur ON r.id = ur.role_id
	WHERE ur.user_id = $1
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RoleRepo) GetRoles(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role

	query := `
	SELECT id, name, description
	FROM auth.roles
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		if errors.As(err, pgx.ErrNoRows) {
			return roles, nil
		} else {
			return nil, err
		}
	}

	for rows.Next() {
		var role models.Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}
