package repository

import (
	"context"
	"core/internal/models"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AudienceRepo struct {
	db *pgxpool.Pool
}

func NewAudienceRepo(db *pgxpool.Pool) *AudienceRepo {
	return &AudienceRepo{db: db}
}

func (r *AudienceRepo) Create(ctx context.Context, audience *models.Audience) error {
	query := `
	INSERT INTO schedule.audiences(name, number)
	VALUES ($1, $2)
	RETURNING id
	`
	return r.db.QueryRow(ctx, query, audience.Name, audience.Number).Scan(&audience.ID)
}

func (r *AudienceRepo) Get(ctx context.Context, id int) (*models.Audience, error) {
	var audience models.Audience
	query := `
	SELECT id, name, number
	FROM schedule.audiences
	WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&audience.ID,
		&audience.Name,
		&audience.Number,
	)
	if err != nil {
		return nil, err
	}

	return &audience, nil
}

func (r *AudienceRepo) GetByNumber(ctx context.Context, number string) (*models.Audience, error) {
	var audience models.Audience
	query := `
	SELECT id, name, number
	FROM schedule.audiences
	WHERE number = $1
	`
	err := r.db.QueryRow(ctx, query, number).Scan(
		&audience.ID,
		&audience.Name,
		&audience.Number,
	)
	if err != nil {
		return nil, err
	}

	return &audience, nil

}

func (r *AudienceRepo) GetAll(ctx context.Context) ([]models.Audience, error) {
	audiences := []models.Audience{}
	query := `
	SELECT id, name, number
	FROM schedule.audiences
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return audiences, nil
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var audience models.Audience
		err := rows.Scan(
			&audience.ID,
			&audience.Name,
			&audience.Number,
		)
		if err != nil {
			return nil, err
		}

		audiences = append(audiences, audience)
	}

	return audiences, nil
}

func (r *AudienceRepo) Update(ctx context.Context, audience *models.Audience) error {
	query := `
	UPDATE schedule.audiences
	SET name = $1, number = $2
	WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, audience.Name, audience.Number, audience.ID)
	return err
}

func (r *AudienceRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM schedule.audiences WHERE id = $1`, id)
	return err
}
