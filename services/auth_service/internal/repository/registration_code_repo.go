package repository

import (
	"auth_service/internal/config"
	"auth_service/internal/dto"
	"auth_service/internal/models"
	"database/sql"
	"errors"
	"math/rand"
	"time"
)

type RegCodeRepository struct {
	db *sql.DB
	cfg *config.Config
}

func NewCodeRepository(db *sql.DB, cfg *config.Config) *RegCodeRepository {
	return &RegCodeRepository{db: db, cfg: cfg}
}

func (r *RegCodeRepository) generateUniqueCode() string {
	charset := r.cfg.CodeCharset
	codeLenght := r.cfg.CodeLength

	code := make([]byte, codeLenght)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func (r *RegCodeRepository) CheckCodeExists(code string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM registration_codes WHERE code = $1)"
	var exists bool
	err := r.db.QueryRow(query, code).Scan(&exists)
	return exists, err
}

func (r *RegCodeRepository) CreateCode(req *dto.CreateRegistrationCodeRequest) (*models.RegistrationCode, error) {
	var role models.Role
	err := r.db.QueryRow("SELECT id FROM roles WHERE name = $1", req.RoleName).Scan(&role.ID)
	if err != nil {
		return nil, err
	}

	maxTries := r.cfg.MaxGenerateTries
	var code string
	i := 0
	for i < maxTries {
		i++
		code = r.generateUniqueCode()

		exists, err := r.CheckCodeExists(code)
		if err != nil {
			return nil, err
		}

		if !exists {
			break
		}

		if i == maxTries {
			return nil, errors.New("failed to generate unique code")
		}
	}

	query := `
	INSERT INTO registration_codes (code, role_id, group_name, max_uses, created_by, expires_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, code, role_id, group_name, max_uses, current_uses, created_by, expires_at, created_at
	`

	expiriesAt := time.Now().Add(time.Duration(req.Expiration) * time.Second)

	var registrationCode models.RegistrationCode
	err = r.db.QueryRow(
		query,
		code,
		role.ID,
		req.GroupName,
		req.MaxUses,
		req.CreatedBy,
		expiriesAt,
	).Scan(
		&registrationCode.ID,
		&registrationCode.Code,
		&registrationCode.RoleID,
		&registrationCode.GroupName,
		&registrationCode.MaxUses,
		&registrationCode.CurrentUses,
		&registrationCode.CreatedBy,
		&registrationCode.ExpiresAt,
		&registrationCode.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	r.db.QueryRow(
		"SELECT id, name, description FROM roles WHERE id = $1",
		registrationCode.RoleID,
	).Scan(
		&registrationCode.Role.ID,
		&registrationCode.Role.Name,
		&registrationCode.Role.Description,
	)


	return &registrationCode, nil
}

func (r *RegCodeRepository) GetCode(code string) (*models.RegistrationCode, error) {
	query := `
	SELECT id, code, role_id, group_name, max_uses, current_uses, created_by, expires_at, created_at
	FROM registration_codes
	WHERE code = $1
	`
	var registrationCode models.RegistrationCode
	err := r.db.QueryRow(
		query,
		code,
	).Scan(
		&registrationCode.ID,
		&registrationCode.Code,
		&registrationCode.RoleID,
		&registrationCode.GroupName,
		&registrationCode.MaxUses,
		&registrationCode.CurrentUses,
		&registrationCode.CreatedBy,
		&registrationCode.ExpiresAt,
		&registrationCode.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	r.db.QueryRow(
		"SELECT id, name, description FROM roles WHERE id = $1",
		registrationCode.RoleID,
	).Scan(
		&registrationCode.Role.ID,
		&registrationCode.Role.Name,
		&registrationCode.Role.Description,
	)

	if registrationCode.GroupName == nil {
		nullGroup := ""
		registrationCode.GroupName = &nullGroup
	}
	
	return &registrationCode, nil
}

func (r *RegCodeRepository) UseCode(code string) error {
	query := `
	UPDATE registration_codes
	SET current_uses = current_uses + 1
	WHERE code = $1
	AND expires_at > NOW()
	AND current_uses < max_uses
	`

	_, err := r.db.Exec(query, code)
	return err
}
