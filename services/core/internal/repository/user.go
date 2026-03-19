package repository

import (
	"context"
	"core/internal/dto"
	"core/internal/models"
	"fmt"
	"math"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db                *pgxpool.Pool
	allowedSortFields map[string]bool
	allowedSortOrders map[string]bool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
		allowedSortFields: map[string]bool{
			"id":         true,
			"username":   true,
			"full_name":  true,
			"email":      true,
			"created_at": true,
		},
		allowedSortOrders: map[string]bool{
			"ASC":  true,
			"DESC": true,
		},
	}
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
	INSERT INTO auth.users (username, full_name, email, password_hash) 
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, query,
		user.Name,
		user.FullName,
		user.Email,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	query = `
	INSERT INTO auth.user_roles (user_id, role_id)
	VALUES (
		$1,
		(SELECT id FROM auth.roles WHERE name = 'user' LIMIT 1)
	)
	ON CONFLICT (user_id, role_id) DO NOTHING
	`
	_, err = tx.Exec(ctx, query, user.ID)
	if err != nil {
		return err
	}

	query = `
	SELECT r.id, r.name, r.description
	FROM auth.user_roles ur
	JOIN auth.roles r ON ur.role_id = r.id
	WHERE ur.user_id = $1
	`

	row, err := tx.Query(ctx, query, user.ID)
	if err == nil {
		for row.Next() {
			var role models.Role
			err := row.Scan(&role.ID, &role.Name, &role.Description)
			if err == nil {
				user.Roles = append(user.Roles, role)
			}
		}
	}
	return tx.Commit(ctx)
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	query := `
	SELECT id, username, full_name, email, password_hash, created_at, updated_at
	FROM auth.users
	WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	query = `
	SELECT r.id, r.name, r.description
	FROM auth.roles r
	JOIN auth.user_roles ur ON ur.role_id = r.id
	WHERE ur.user_id = $1
	`
	row, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var role models.Role
		err = row.Scan(&role.ID, &role.Name, &role.Description)
		if err != nil {
			return nil, err
		}
		user.Roles = append(user.Roles, role)
	}

	return &user, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	query := `
	SELECT id, username, full_name, email, password_hash, created_at, updated_at
	FROM auth.users
	WHERE username = $1
	`
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Name,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	query = `
	SELECT r.id, r.name, r.description
	FROM auth.roles r
	JOIN auth.user_roles ur ON ur.role_id = r.id
	WHERE ur.user_id = $1
	`
	row, err := r.db.Query(ctx, query, user.ID)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var role models.Role
		err = row.Scan(&role.ID, &role.Name, &role.Description)
		if err != nil {
			return nil, err
		}
		user.Roles = append(user.Roles, role)
	}

	return &user, nil
}

func (r *UserRepo) GetUsersPaginated(ctx context.Context, page, perPage int, sortBy, sortOrder string) (*dto.PaginatedUsers, error) {
	if page < 1 {
		page = 1
	}

	if perPage < 10 {
		page = 10
	}

	if sortBy == "" || !r.allowedSortFields[sortBy] {
		sortBy = "id"
	}

	sortOrder = strings.ToUpper(sortOrder)
	if sortOrder == "" || !r.allowedSortOrders[sortOrder] {
		sortOrder = "ASC"
	}

	offset := (page - 1) * perPage
	orderClause := fmt.Sprintf("%s %s", sortBy, sortOrder)

	query := fmt.Sprintf(`
	SELECT id, username, full_name, email, created_at, updated_at
	FROM auth.users
	ORDER BY %s
	LIMIT $1 OFFSET $2
	`, orderClause)

	rows, err := r.db.Query(ctx, query, perPage, page)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.FullName,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}

		users = append(users, user)
	}

	var total int
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM auth.users`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count users: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.PaginatedUsers{
		Users: users,
		Total: total,
		Page: page,
		PerPage: perPage,
		TotalPages: totalPages,
	}, nil
}

func (r *UserRepo) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User

	query := `
	SELECT id, username, full_name, email, created_at, updated_at
	FROM auth.users 
	`

	row, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		var user models.User
		err := row.Scan(
			&user.ID,
			&user.Name,
			&user.FullName,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepo) UpdatePasswordHash(ctx context.Context, id int, passwordHash string) error {
	query := `
	UPDATE auth.users
	SET password_hash = $1
	WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, passwordHash, id)
	return err
}

func (r *UserRepo) UpdateFullName(ctx context.Context, id int, fullName string) error {
	query := `
	UPDATE auth.users
	SET full_name = $1
	WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, fullName, id)
	return err
}

func (r *UserRepo) Delete(ctx context.Context, id int) error {
	query := `
	DELETE FROM auth.users WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
