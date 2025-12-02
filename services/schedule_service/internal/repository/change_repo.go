package repository

import (
	"schedule-service/internal/models"
	"time"
)

type ChangeRepository interface {
	AddChange(sc *models.ScheduleChange) error
	GetChange(date time.Time) ([]models.ScheduleChange, error)
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
	INSERT INTO changes(change_date, change_image_url, change_description)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`

	err = tx.QueryRow(query, sc.Date, sc.ImgURL, sc.Description).Scan(&sc.ID, &sc.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *changeRepo) GetChange(date time.Time) ([]models.ScheduleChange, error) {
	query := `
	SELECT id, change_date, change_image_url, change_description, created_at
	FROM changes
	WHERE change_date = $1
	`
	chs := []models.ScheduleChange{}
	rows, err := r.db.DB.Query(query, date)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		ch := models.ScheduleChange{}
		if err := rows.Scan(
			&ch.ID,
			&ch.Date,
			&ch.ImgURL,
			&ch.Description,
			&ch.CreatedAt,
		); err != nil {
			return nil, err
		}

		chs = append(chs, ch)
	}

	return chs, nil
}

func (r *changeRepo) RemoveChange(id int) error {
	query := `
	DELETE FROM changes
	WHERE id = $1
	`

	_, err := r.db.DB.Exec(query, id)
	return err
}
