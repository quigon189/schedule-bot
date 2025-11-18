package repository

import (
	"auth_service/internal/models"
	"fmt"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	Get(telegramID int64) (*models.User, error)
	Update(user *models.User) error
	UpdateUserRoles(telegramID int64, roles []string) error
	Delete(telegramID int64) error
	List() ([]models.User, error)
}

type userRepo struct {
	db *DB
}

func NewUserRepository(db *DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *models.User) (*models.User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO users (telegram_id, username, full_name)
	VALUES ($1, $2, $3) RETURNING created_at, updated_at, is_active
	`
	err = tx.QueryRow(query, user.TelegramID, user.Username, user.FullName).
		Scan(&user.CreatedAt, &user.UpdatedAt, &user.IsActive)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	var role models.Role
	err = tx.QueryRow("SELECT id, name, description FROM roles WHERE name = $1", "user").
		Scan(&role.ID, &role.Name, &role.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to find role user: %v", err)
	}

	_, err = tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)", user.TelegramID, role.ID)
	if err != nil {
		return nil, fmt.Errorf("faile to assign role: %v", err)
	}

	user.Roles = append(user.Roles, role)

	return user, tx.Commit()
}

func (r *userRepo) Get(telegramID int64) (*models.User, error) {
	user := &models.User{}
	query := `
	SELECT telegram_id, username, full_name, created_at, updated_at, is_active
	FROM users WHERE telegram_id = $1
	`
	err := r.db.QueryRow(query, telegramID).Scan(
		&user.TelegramID, &user.Username, &user.FullName,
		&user.CreatedAt, &user.UpdatedAt, &user.IsActive,
	)
	if err != nil {
		return nil, err
	}

	query = `
	SELECT r.id, r.name, r.description
	FROM roles r
	JOIN user_roles ur ON r.id = ur.role_id
	WHERE ur.user_id = $1
	`
	rows, err := r.db.Query(query, user.TelegramID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description)
		if err != nil {
			return nil, err
		}

		user.Roles = append(user.Roles, role)
	}

	var student models.Student
	err = r.db.QueryRow("SELECT user_id, group_name FROM students WHERE user_id = $1", user.TelegramID).
		Scan(&student.UserID, &student.GroupName)
	if err == nil {
		user.Student = &student
	}

	return user, nil
}

func (r *userRepo) Update(user *models.User) error {
	dbUser, err := r.Get(user.TelegramID)
	if err != nil {
		return fmt.Errorf("failed to get user %d: %v", user.TelegramID, err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if dbUser.Username != user.Username || dbUser.FullName != user.FullName || dbUser.IsActive != user.IsActive {
		query := "UPDATE users SET updated_at = CURRENT_TIMESTAMP"
		params := []any{}
		paramCount := 1

		if dbUser.Username != user.Username {
			query += fmt.Sprintf(", username = $%d", paramCount)
			params = append(params, user.Username)
			paramCount++
		}
		if dbUser.FullName != user.FullName {
			query += fmt.Sprintf(", full_name = $%d", paramCount)
			params = append(params, user.FullName)
			paramCount++
		}
		if dbUser.IsActive != user.IsActive {
			query += fmt.Sprintf(", is_active = $%d", paramCount)
			params = append(params, user.IsActive)
			paramCount++
		}

		query += fmt.Sprintf(" WHERE telegram_id = $%d", paramCount)
		params = append(params, user.TelegramID)

		_, err := tx.Exec(query, params...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *userRepo) UpdateUserRoles(telegramID int64, roles []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM user_roles WHERE user_id = $1", telegramID)
	if err != nil {
		return err
	}

	for _, role := range roles {
		var roleID int
		err := tx.QueryRow("SELECT id FROM roles WHERE name = $1", role).Scan(&roleID)
		if err != nil {
			return err
		}

		_, err = tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)", telegramID, roleID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *userRepo) Delete(telegramID int64) error {
	_, err := r.db.Exec("DELETE FROM users WHERE telegram_id = $1", telegramID)
	return err
}

func (r *userRepo) List() ([]models.User, error) {
	users := []models.User{}
	query := `
	SELECT telegram_id, username, full_name, created_at, updated_at, is_active
	FROM users
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user := models.User{}
		err := rows.Scan(&user.TelegramID, &user.Username, &user.FullName,
			&user.CreatedAt, &user.UpdatedAt, &user.IsActive)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
