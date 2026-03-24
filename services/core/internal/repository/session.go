package repository

import (
	"context"
	"core/internal/models"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepo struct {
	db *pgxpool.Pool
}

func NewSessionRepo(db *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{db: db}
}

func (r *SessionRepo) Create(ctx context.Context, session *models.Session) error {
	query := `
	INSERT INTO auth.sessions (user_id, refresh_token, user_agent, client_ip)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query,
		session.User.ID,
		session.RefreshToken,
		session.UserAgent,
		session.ClientIP,
	).Scan(
		&session.ID,
		&session.CreatedAt,
	)
}

func (r *SessionRepo) GetByID(ctx context.Context, id string) (*models.Session, error) {
	session := models.Session{
		User: &models.User{},
	}

	query := `
	SELECT s.id, s.refresh_token, s.user_agent, s.client_ip,
		u.id, u.username, u.full_name, u.email, u.password_hash,
		u.created_at, u.updated_at
	FROM auth.sessions s
	JOIN auth.users u ON u.id = s.user_id
	WHERE s.id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&session.ID,
		&session.RefreshToken,
		&session.UserAgent,
		&session.ClientIP,
		&session.User.ID,
		&session.User.Name,
		&session.User.FullName,
		&session.User.Email,
		&session.User.PasswordHash,
		&session.User.CreatedAt,
		&session.User.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepo) GetByUserID(ctx context.Context, id int) ([]models.Session, error) {
	sessions := []models.Session{}
	query := `
	SELECT id, refresh_token, user_agent, client_ip
	FROM auth.sessions
	WHERE user_id = $1
	`
	row, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		session := models.Session{}
		err := row.Scan(
			&session.ID,
			&session.RefreshToken,
			&session.UserAgent,
			&session.ClientIP,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *SessionRepo) UpdateRefreshToken(ctx context.Context, id string, refreshToken string) error {
	query := `
	UPDATE auth.sessions
	SET refresh_token = $1
	WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, refreshToken, id)
	return err
}

func (r *SessionRepo) Delete(ctx context.Context, id pgtype.UUID) error {
	query := `
	DELETE FROM auth.sessions WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
