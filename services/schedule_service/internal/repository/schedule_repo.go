package repository

import (
	"database/sql"
	"fmt"
	"schedule-service/internal/models"
)

type ScheduleRepository interface {
	AddGroupSchedule(gs *models.GroupSchedule) error
	// UpdateGroupSchedule(sp *models.StudyPeriod, scheduleImgURL string) (*models.GroupSchedule, error)
	GetGroupSchedules(filter *models.GroupScheduleFilter) ([]models.GroupSchedule, error)
	RemoveGroupSchedule(gs *models.GroupSchedule) error
}

type scheduleRepo struct {
	db *DB
}

func NewScheduleRepository(db *DB) ScheduleRepository {
	return &scheduleRepo{db: db}
}

func (r *scheduleRepo) AddGroupSchedule(gs *models.GroupSchedule) error {
	tx, err := r.db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = getStudyPeriod(tx, &gs.StudyPeriod); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Not found study period")
		}
		return fmt.Errorf("Failed get study period: %v", err)
	}

	if err = getGroup(tx, &gs.Group); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Not found group")
		}
		return fmt.Errorf("Failed get group: %v", err)
	}

	query := `
	INSERT INTO group_schedules (study_period_id, group_id, semester, schedule_image_url)
	VALUES ($1, $2, $3, $4)
	RETURNING created_at
	`
	err = tx.QueryRow(query, gs.StudyPeriod.ID, gs.Group.ID, gs.Semester, gs.ScheduleImgURL).Scan(&gs.CreatedAt)
	if err != nil {
		return fmt.Errorf("Failed to get group schedule create time: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// func (r *scheduleRepo) UpdateGroupSchedule(sp *models.StudyPeriod, scheduleImgURL string) (*models.GroupSchedule, error) {
// 	return &models.GroupSchedule{}, nil
// }

func (r *scheduleRepo) GetGroupSchedules(filter *models.GroupScheduleFilter) ([]models.GroupSchedule, error) {
	tx, err := r.db.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	SELECT study_period_id, academic_year, half_year, start_date, end_date,
		   group_id, group_name, semester, schedule_image_url, created_at
	FROM group_schedules_view
	`
	paramsCount := 0
	params := []any{}

	if filter.HalfYear != 0 ||
		filter.AcademicYear != "" ||
		filter.GroupName != "" ||
		filter.Semester != 0 {
		query += "WHERE"

		if filter.AcademicYear != "" {
			paramsCount++
			params = append(params, filter.AcademicYear)
			query += fmt.Sprintf(" academic_year = $%d", paramsCount)
		}
		if filter.HalfYear != 0 {
			if paramsCount != 0 {
				query += " AND"
			}
			paramsCount++
			params = append(params, filter.HalfYear)
			query += fmt.Sprintf(" half_year = $%d", paramsCount)
		}
		if filter.Semester != 0 {
			if paramsCount != 0 {
				query += " AND"
			}
			paramsCount++
			params = append(params, filter.Semester)
			query += fmt.Sprintf(" semester = $%d", paramsCount)
		}
		if filter.GroupName != "" {
			if paramsCount != 0 {
				query += " AND"
			}
			paramsCount++
			params = append(params, filter.GroupName)
			query += fmt.Sprintf(" group_name = $%d", paramsCount)
		}
	}

	gss := []models.GroupSchedule{}

	rows, err := tx.Query(query, params...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		gs := models.GroupSchedule{}
		err := rows.Scan(
			&gs.StudyPeriod.ID,
			&gs.StudyPeriod.AcademicYear,
			&gs.StudyPeriod.HalfYear,
			&gs.StudyPeriod.StrartDate,
			&gs.StudyPeriod.EndDate,
			&gs.Group.ID,
			&gs.Group.Name,
			&gs.Semester,
			&gs.ScheduleImgURL,
			&gs.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		gss = append(gss, gs)
	}

	return gss, nil
}

func (r *scheduleRepo) RemoveGroupSchedule(gs *models.GroupSchedule) error {
	query := `
	DELETE FROM group_schedules
	WHERE study_period_id = $1 AND group_id = $2
	`
	_, err := r.db.DB.Exec(query, gs.StudyPeriod.ID, gs.Group.ID)
	return err
}

func getGroup(tx *sql.Tx, g *models.Group) error {
	query := `
	SELECT id, name, created_at
	FROM groups
	WHERE name = $1
	`
	return tx.QueryRow(query, g.Name).Scan(&g.ID, &g.Name, &g.CreatedAt)
}

func getStudyPeriod(tx *sql.Tx, sp *models.StudyPeriod) error {
	query := `
	SELECT id, half_year, academic_year, start_date, end_date, created_at
	FROM study_periods
	WHERE half_year = $1 AND academic_year = $2
	`
	return tx.QueryRow(
		query,
		sp.HalfYear,
		sp.AcademicYear,
	).Scan(
		&sp.ID,
		&sp.HalfYear,
		&sp.AcademicYear,
		&sp.StrartDate,
		&sp.EndDate,
		&sp.CreatedAt,
	)
}
