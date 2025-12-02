package service

import (
	"fmt"
	"schedule-service/internal/dto"
	"schedule-service/internal/models"
	"schedule-service/internal/repository"
)

type ScheduleService struct {
	repo repository.ScheduleRepository
}

func NewScheduleService(repo *repository.ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: *repo}
}

func (s *ScheduleService) AddGroupSchedule(req *dto.AddGroupScheduleRequest) (*models.GroupSchedule, error) {
	gs := &models.GroupSchedule{
		Group: models.Group{
			Name: req.GroupName,
		},
		StudyPeriod: models.StudyPeriod{
			HalfYear:     req.HalfYear,
			AcademicYear: req.AcademicYear,
		},
		Semester:       req.Semester,
		ScheduleImgURL: req.ScheduleImgURL,
	}

	err := s.repo.AddGroupSchedule(gs)
	if err != nil {
		return nil, fmt.Errorf("Failed to create group schedule: %v", err)
	}

	return gs, nil
}

func (s *ScheduleService) GetGroupSchedule(req *dto.GroupScheduleQueryParams) ([]models.GroupSchedule, error) {
	filter := &models.GroupScheduleFilter{
		AcademicYear: req.AcademicYear,
		HalfYear:     req.HalfYear,
		GroupName:    req.GroupName,
	}

	gs, err := s.repo.GetGroupSchedules(filter)
	if err != nil {
		return nil, fmt.Errorf("Failed to get group schedule: %v", err)
	}

	return gs, err
}

func (s *ScheduleService) RemoveScheduleService(req *dto.GroupScheduleQueryParams) error {
	filter := &models.GroupScheduleFilter{
		AcademicYear: req.AcademicYear,
		HalfYear:     req.HalfYear,
		GroupName:    req.GroupName,
	}
	gs, err := s.repo.GetGroupSchedules(filter)
	if err != nil {
		return err
	}

	err = s.repo.RemoveGroupSchedule(&gs[0])
	if err != nil {
		return fmt.Errorf("Failed to remove group schedule: %v", err)
	}

	return nil
}
