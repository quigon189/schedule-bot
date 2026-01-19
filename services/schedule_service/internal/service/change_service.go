package service

import (
	"fmt"
	"schedule-service/internal/dto"
	"schedule-service/internal/models"
	"schedule-service/internal/repository"
	"time"
)

type ChangeService struct {
	repo repository.ChangeRepository
}

func NewChangeService(repo repository.ChangeRepository) *ChangeService {
	return &ChangeService{repo: repo}
}

func (s *ChangeService) AddChange(req dto.AddChangeRequest) (*models.ScheduleChange, error) {
	date, err := parseDate(req.Date)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse date: %v", err)
	}
	sc := models.ScheduleChange{
		Date:        date,
		ImgURLs:     req.ImgURLs,
		Description: req.Description,
	}

	if err := s.repo.AddChange(&sc); err != nil {
		return nil, fmt.Errorf("Failed to adding change: %v", err)
	}

	return &sc, nil
}

func (s *ChangeService) GetChanges(dateStr string) (*models.ScheduleChange, error) {
	var date time.Time
	var err error
	if dateStr == "" {
		date = time.Now()
	} else {
		date, err = parseDate(dateStr)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse date: %v", err)
		}
	}

	return s.repo.GetChange(date)
}

func (s *ChangeService) RemoveChange(id int) error {
	return s.repo.RemoveChange(id)
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}
