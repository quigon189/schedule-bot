package services

import (
	"context"
	"core/internal/models"
	"core/internal/repository"
)

type ScheduleService struct {
	audienceRepo *repository.AudienceRepo
}

func NewScheduleService(audienceRepo *repository.AudienceRepo) *ScheduleService {
	return &ScheduleService{audienceRepo: audienceRepo}
}

func (s *ScheduleService) CreateAudience(ctx context.Context, audience *models.Audience) error {
	return s.audienceRepo.Create(ctx, audience)
}

func (s *ScheduleService) GetAudience(ctx context.Context, id int) (*models.Audience, error) {
	return s.audienceRepo.Get(ctx, id)
}

func (s *ScheduleService) GetAllAudience(ctx context.Context) ([]models.Audience, error) {
	return s.audienceRepo.GetAll(ctx)
}

func (s *ScheduleService) UpdateAudience(ctx context.Context, audience *models.Audience) error {
	return s.audienceRepo.Update(ctx, audience)
}

func (s *ScheduleService) DeleteAudience(ctx context.Context, id int) error {
	return s.audienceRepo.Delete(ctx, id)
}
