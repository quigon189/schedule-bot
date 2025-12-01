package repository

import (
	"fmt"
	"schedule-service/internal/models"
	"time"
)

type ChangeRepository interface{
	AddChange(sc *models.ScheduleChange) error
	GetChange(date time.Time) (*models.Change, error)
	RemoveChange(sc *models.Change) error
}

type changeRepo struct {
	db *DB
}

func NewChangeRepo(db *DB) *ChangeRepository {
	return &changeRepo{db: db}
}

func (r *changeRepo) AddChange(sc *models.Change) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	SELECT id FROM changes WHERE change_image_url = $1
	`
	_, err = tx.Query(query, sc.ImgURL)
	if err == nil {
		return fmt.Errorf("duplicate image url")
	}

	query = `
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

func (r *changeRepo) GetChange(date *time.Time) (*models.Change, error) {
	query := `
	SELECT id, change_date, change_image_url, change_description, created_at
	FROM changes
	WHERE change_date = $1
	`
	rows, err := r.db.DB.Query(query, date)
	if err != nil{
		return nil, err
	}

}
