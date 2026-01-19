package repository

import (
	"github.com/lib/pq"
	"schedule-service/internal/models"
	"time"
)

type ChangeRepository interface {
	AddChange(sc *models.ScheduleChange) error
	GetChange(date time.Time) (*models.ScheduleChange, error)
	RemoveChange(id int) error
}

type changeRepo struct {
	db *DB
}

func NewChangeRepo(db *DB) ChangeRepository {
	return &changeRepo{db: db}
}

func (r *changeRepo) AddChange(sc *models.ScheduleChange) error {
	tx, err := r.db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO changes(change_date, change_image_urls, change_description)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`

	err = tx.QueryRow(query, sc.Date, pq.Array(sc.ImgURLs), sc.Description).Scan(&sc.ID, &sc.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *changeRepo) GetChange(date time.Time) (*models.ScheduleChange, error) {
	query := `
	SELECT id, change_date, change_image_urls, change_description, created_at
	FROM changes
	WHERE change_date = $1
	`
	ch := models.ScheduleChange{}
	err := r.db.DB.QueryRow(query, date).Scan(
		&ch.ID,
		&ch.Date,
		pq.Array(&ch.ImgURLs),
		&ch.Description,
		&ch.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &ch, nil
}

func (r *changeRepo) RemoveChange(id int) error {
	query := `
	DELETE FROM changes
	WHERE id = $1
	`

	_, err := r.db.DB.Exec(query, id)
	return err
}
