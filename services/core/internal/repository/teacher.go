package repository

import (
	"context"
	"core/internal/models"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeacherRepo struct {
	db *pgxpool.Pool
}

func NewTeacherRepo(db *pgxpool.Pool) *TeacherRepo {
	return &TeacherRepo{db: db}
}

func (r *TeacherRepo) CreateTeacher(ctx context.Context, user *models.User) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
	INSERT INTO auth.user_roles (user_id, role_id)
	VALUES ($1, (SELECT id FROM auth.roles WHERE name = 'teacher'))
	ON CONFLICT DO NOTHING
	`
	if _, err = tx.Exec(ctx, query, user.ID); err != nil {
		return err
	}

	query = `
	INSERT INTO auth.teacher_profiles (user_id)
	VALUES ($1)
	`
	if _, err = tx.Exec(ctx, query, user.ID); err != nil {
		return err
	}

	query = `
	SELECT r.id, r.name, r.description
	FROM auth.user_roles ur
	JOIN auth.roles r ON ur.role_id = r.id
	WHERE ur.user_id = $1
	`
	rows, err := tx.Query(ctx, query, user.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description); err != nil {
			return err
		}
		user.Roles = append(user.Roles, role)
	}

	return tx.Commit(ctx)
}

func (r *TeacherRepo) GetTeacherByUserID(ctx context.Context, userID int) (*models.Teacher, error) {
	var teacher models.Teacher
	query := `
	SELECT u.id, u.username, u.full_name, u.email, u.password_hash, u.created_at, u.updated_at
	FROM auth.users u
	JOIN auth.teacher_profiles tp ON u.id = tp.user_id
	WHERE u.id = $1
	`
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&teacher.User.ID,
		&teacher.User.Name,
		&teacher.User.FullName,
		&teacher.User.Email,
		&teacher.User.PasswordHash,
		&teacher.User.CreatedAt,
		&teacher.User.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	query = `
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
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description); err != nil {
			return nil, err
		}
		teacher.User.Roles = append(teacher.User.Roles, role)
	}

	return &teacher, nil
}

func (r *TeacherRepo) GetAllTeachers(ctx context.Context) ([]models.Teacher, error) {
	query := `
	SELECT u.id, u.username, u.full_name, u.email, u.password_hash, u.created_at, u.updated_at
	FROM auth.users u
	JOIN auth.teacher_profiles tp ON u.id = tp.user_id
	ORDER BY u.id
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []models.Teacher
	for rows.Next() {
		var t models.Teacher
		if err := rows.Scan(
			&t.User.ID,
			&t.User.Name,
			&t.User.FullName,
			&t.User.Email,
			&t.User.PasswordHash,
			&t.User.CreatedAt,
			&t.User.UpdatedAt,
		); err != nil {
			return nil, err
		}
		teachers = append(teachers, t)
	}
	return teachers, nil
}

func (r *TeacherRepo) DeleteTeacher(ctx context.Context, userID int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM auth.teacher_profiles WHERE user_id = $1`, userID)
	return err
}

func ( r *TeacherRepo) IsTeacher(ctx context.Context, userID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM auth.teacher_profiles WHERE user_id = $1)`
	err := r.db.QueryRow(ctx, query, userID).Scan(&exists)
	return exists, err
}
