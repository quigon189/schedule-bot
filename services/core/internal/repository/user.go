package repository

import (
	"context"
	"core/internal/dto"
	"core/internal/models"
	"fmt"
	"math"
	"strings"

	"github.com/jackc/pgx/v5"
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

	if err := r.CreateWithTx(ctx, tx, user); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *UserRepo) CreateWithTx(ctx context.Context, tx pgx.Tx, user *models.User) error {

	query := `
	INSERT INTO auth.users (username, full_name, email, password_hash) 
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at
	`
	err := tx.QueryRow(ctx, query,
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
	row.Close()

	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	query := `
	SELECT u.id, u.username, u.full_name, u.email, u.created_at, u.updated_at, u.password_hash,
	COALESCE(
        (SELECT json_agg(json_build_object('id', r.id, 'name', r.name, "description", r.description))
         FROM auth.roles r
         JOIN auth.user_roles ur ON ur.role_id = r.id
         WHERE ur.user_id = u.id), 
        '[]'
    ) as roles,
    (SELECT json_build_object('id', g.id, 'name', g.name, 'specialty', g.specialty, 'admission_year', g.admission_year)
     FROM auth.groups g
     JOIN auth.student_profiles s ON s.group_id = g.id
     WHERE s.user_id = u.id
     -- Проверка на наличие роли student (опционально, если логика БД это гарантирует)
     AND EXISTS (
         SELECT 1 FROM auth.roles r 
         JOIN auth.user_roles ur ON ur.role_id = r.id 
         WHERE ur.user_id = u.id AND r.name = 'student'
     )
     LIMIT 1
    ) as "group"
	FROM auth.users u
	WHERE id = $1
	`
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	user, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	query := `
	SELECT u.id, u.username, u.full_name, u.email, u.created_at, u.updated_at, u.password_hash,
	COALESCE(
        (SELECT json_agg(json_build_object('id', r.id, 'name', r.name, "description", r.description))
         FROM auth.roles r
         JOIN auth.user_roles ur ON ur.role_id = r.id
         WHERE ur.user_id = u.id), 
        '[]'
    ) as roles,
    (SELECT json_build_object('id', g.id, 'name', g.name, 'specialty', g.specialty, 'admission_year', g.admission_year)
     FROM auth.groups g
     JOIN auth.student_profiles s ON s.group_id = g.id
     WHERE s.user_id = u.id
     -- Проверка на наличие роли student (опционально, если логика БД это гарантирует)
     AND EXISTS (
         SELECT 1 FROM auth.roles r 
         JOIN auth.user_roles ur ON ur.role_id = r.id 
         WHERE ur.user_id = u.id AND r.name = 'student'
     )
     LIMIT 1
    ) as "group"
	FROM auth.users u
	WHERE username = $1
	`

	rows, err := r.db.Query(ctx, query, username)
	if err != nil {
		return nil, err
	}
	user, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetUsersPaginated(ctx context.Context, filters *dto.UserFilter, page, perPage int, sortBy, sortOrder string) (*dto.PaginatedUsers, error) {
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

	query := `
	SELECT u.id, u.username, u.full_name, u.email, u.created_at, u.updated_at, u.password_hash,
	COALESCE(
        (SELECT json_agg(json_build_object('id', r.id, 'name', r.name, "description", r.description))
         FROM auth.roles r
         JOIN auth.user_roles ur ON ur.role_id = r.id
         WHERE ur.user_id = u.id), 
        '[]'
    ) as roles,
    (SELECT json_build_object('id', g.id, 'name', g.name, 'specialty', g.specialty, 'admission_year', g.admission_year)
     FROM auth.groups g
     JOIN auth.student_profiles s ON s.group_id = g.id
     WHERE s.user_id = u.id
     -- Проверка на наличие роли student (опционально, если логика БД это гарантирует)
     AND EXISTS (
         SELECT 1 FROM auth.roles r 
         JOIN auth.user_roles ur ON ur.role_id = r.id 
         WHERE ur.user_id = u.id AND r.name = 'student'
     )
     LIMIT 1
    ) as "group"
	FROM auth.users u`

	var conditions []string
	var args []any
	argIndex := 1

	if filters.FullName != nil {
		conditions = append(conditions, fmt.Sprintf("full_name ILIKE '%%' || $%d || '%%'", argIndex))
		args = append(args, *filters.FullName)
		argIndex++
	}

	if filters.Username != nil {
		conditions = append(conditions, fmt.Sprintf("username ILIKE '%%' || $%d || '%%'", argIndex))
		args = append(args, *filters.Username)
		argIndex++
	}

	if filters.Email != nil {
		conditions = append(conditions, fmt.Sprintf("email ILIKE '%%' || $%d || '%%'", argIndex))
		args = append(args, *filters.Email)
		argIndex++
	}

	if filters.GroupName != nil {
		conditions = append(conditions,
			fmt.Sprintf(`EXISTS (
            SELECT 1 FROM auth.student_profiles s
            JOIN auth.groups g ON g.id = s.group_id
            WHERE s.user_id = u.id
              AND g.name ILIKE '%%' || $%d || '%%'
        )`, argIndex))
		args = append(args, *filters.GroupName)
		argIndex++
	}

	if filters.Role != nil {
		conditions = append(conditions,
			fmt.Sprintf(`EXISTS (
            SELECT 1 FROM auth.user_roles ur
            JOIN auth.roles r ON r.id = ur.role_id
            WHERE ur.user_id = u.id
              AND r.name ILIKE '%%' || $%d || '%%'
        )`, argIndex))
		args = append(args, *filters.Role)
		argIndex++
	}

	if len(conditions) > 0 {
		query += "\nWHERE " + strings.Join(conditions, " AND ")
	}

	query = query + fmt.Sprintf(`
	ORDER BY %s
	LIMIT $%d OFFSET $%d`, orderClause, argIndex, argIndex+1)
	args = append(args, perPage)
	args = append(args, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, fmt.Errorf("collect rows: %w")
	}

	query = `SELECT COUNT(*) FROM auth.users u`
	if len(conditions) > 0 {
		query += "\nWHERE " + strings.Join(conditions, " AND ")
	}
	args = args[:len(args)-2]

	var total int
	err = r.db.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count users: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.PaginatedUsers{
		Users:      users,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (r *UserRepo) GetAll(ctx context.Context, filters *dto.UserFilter) ([]models.User, error) {
	var conditions []string
	var args []any
	argIndex := 1

	var users []models.User

	query := `
	SELECT u.id, u.username, u.full_name, u.email, u.created_at, u.updated_at, u.password_hash,
	COALESCE(
        (SELECT json_agg(json_build_object('id', r.id, 'name', r.name, "description", r.description))
         FROM auth.roles r
         JOIN auth.user_roles ur ON ur.role_id = r.id
         WHERE ur.user_id = u.id), 
        '[]'
    ) as roles,
    (SELECT json_build_object('id', g.id, 'name', g.name, 'specialty', g.specialty, 'admission_year', g.admission_year)
     FROM auth.groups g
     JOIN auth.student_profiles s ON s.group_id = g.id
     WHERE s.user_id = u.id
     -- Проверка на наличие роли student (опционально, если логика БД это гарантирует)
     AND EXISTS (
         SELECT 1 FROM auth.roles r 
         JOIN auth.user_roles ur ON ur.role_id = r.id 
         WHERE ur.user_id = u.id AND r.name = 'student'
     )
     LIMIT 1
    ) as "group"
	FROM auth.users u
	`

	if filters.FullName != nil {
		conditions = append(conditions, fmt.Sprintf("full_name ILIKE '%%' || $%d || '%%'", argIndex))
		args = append(args, *filters.FullName)
		argIndex++
	}

	if filters.Username != nil {
		conditions = append(conditions, fmt.Sprintf("username ILIKE '%%' || $%d || '%%'", argIndex))
		args = append(args, *filters.Username)
		argIndex++
	}

	if filters.Email != nil {
		conditions = append(conditions, fmt.Sprintf("email ILIKE '%%' || $%d || '%%'", argIndex))
		args = append(args, *filters.Email)
		argIndex++
	}

	if filters.GroupName != nil {
		conditions = append(conditions,
			fmt.Sprintf(`EXISTS (
            SELECT 1 FROM auth.student_profiles s
            JOIN auth.groups g ON g.id = s.group_id
            WHERE s.user_id = u.id
              AND g.name ILIKE '%%' || $%d || '%%'
        )`, argIndex))
		args = append(args, *filters.GroupName)
		argIndex++
	}

	if filters.Role != nil {
		conditions = append(conditions,
			fmt.Sprintf(`EXISTS (
            SELECT 1 FROM auth.user_roles ur
            JOIN auth.roles r ON r.id = ur.role_id
            WHERE ur.user_id = u.id
              AND r.name ILIKE '%%' || $%d || '%%'
        )`, argIndex))
		args = append(args, *filters.Role)
		argIndex++
	}

	if len(conditions) > 0 {
		query += "\nWHERE " + strings.Join(conditions, " AND ")
	}

	row, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	users, err = pgx.CollectRows(row, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, fmt.Errorf("collect rows: %w", err)
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
