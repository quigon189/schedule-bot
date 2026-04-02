package repository

import (
	"context"
	"core/internal/models"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentRepo struct {
	db *pgxpool.Pool
}

func NewStudentRepo(db *pgxpool.Pool) *StudentRepo {
	return &StudentRepo{db: db}
}

func (r *StudentRepo) CreateStudent(ctx context.Context, user *models.User, groupID int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := r.CreateStudentWithTx(ctx, tx, user, groupID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *StudentRepo) CreateStudentWithTx(ctx context.Context, tx pgx.Tx, user *models.User, groupID int) error {
	query := `
	INSERT INTO auth.user_roles (user_id, role_id)
	VALUES ($1, (SELECT id FROM auth.roles WHERE name = 'student'))
	ON CONFLICT DO NOTHING
	`
	if _, err := tx.Exec(ctx, query, user.ID); err != nil {
		return err
	}

	query = `
	INSERT INTO auth.student_profiles (user_id, group_id)
	VALUES ($1, $2)
	`
	if _, err := tx.Exec(ctx, query, user.ID, groupID); err != nil {
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

	return nil
}

func (r *StudentRepo) GetStudentByUserID(ctx context.Context, userID int) (*models.Student, error) {
	student := models.Student{
		User:  models.User{},
		Group: &models.Group{},
	}
	query := `
	SELECT u.id, u.username, u.full_name, u.email, u.password_hash, u.created_at, u.updated_at,
	       g.id, g.name, g.specialty, g.admission_year
	FROM auth.users u
	JOIN auth.student_profiles sp ON u.id = sp.user_id
	LEFT JOIN auth.groups g ON sp.group_id = g.id
	WHERE u.id = $1
	`
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&student.User.ID,
		&student.User.Name,
		&student.User.FullName,
		&student.User.Email,
		&student.User.PasswordHash,
		&student.User.CreatedAt,
		&student.User.UpdatedAt,
		&student.Group.ID,
		&student.Group.Name,
		&student.Group.Specialty,
		&student.Group.AdmissionYear,
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
		student.User.Roles = append(student.User.Roles, role)
	}

	return &student, nil
}

func (r *StudentRepo) GetAllStudents(ctx context.Context) ([]models.Student, error) {
	query := `
	SELECT u.id, u.username, u.full_name, u.email, u.password_hash, u.created_at, u.updated_at,
	       g.id, g.name, g.specialty, g.admission_year
	FROM auth.users u
	JOIN auth.student_profiles sp ON u.id = sp.user_id
	LEFT JOIN auth.groups g ON sp.group_id = g.id
	ORDER BY u.id
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		s := models.Student{
			User:  models.User{},
			Group: &models.Group{},
		}
		if err := rows.Scan(
			&s.User.ID,
			&s.User.Name,
			&s.User.FullName,
			&s.User.Email,
			&s.User.PasswordHash,
			&s.User.CreatedAt,
			&s.User.UpdatedAt,
			&s.Group.ID,
			&s.Group.Name,
			&s.Group.Specialty,
			&s.Group.AdmissionYear,
		); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}

func (r *StudentRepo) UpdateStudentGroup(ctx context.Context, userID int, groupID int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE auth.student_profiles
		SET group_id = $1
		WHERE user_id = $2
	`, groupID, userID)
	return err
}

func (r *StudentRepo) DeleteStudent(ctx context.Context, userID int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM auth.student_profiles WHERE user_id = $1`, userID)
	return err
}

func (r *StudentRepo) IsStudent(ctx context.Context, userID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM auth.student_profiles WHERE user_id = $1)`
	err := r.db.QueryRow(ctx, query, userID).Scan(&exists)
	return exists, err
}
