package services

import (
	"context"
	"core/internal/dto"
	"core/internal/models"
	"core/internal/repository"
	"fmt"
)

type ScheduleService struct {
	audienceRepo *repository.AudienceRepo
	groupRepo    *repository.GroupRepo
	teacherRepo  *repository.TeacherRepo
	studentRepo  *repository.StudentRepo
}

func NewScheduleService(
	audienceRepo *repository.AudienceRepo,
	groupRepo *repository.GroupRepo,
	teacherRepo *repository.TeacherRepo,
	studentRepo *repository.StudentRepo,
) *ScheduleService {
	return &ScheduleService{
		audienceRepo: audienceRepo,
		groupRepo:    groupRepo,
		teacherRepo:  teacherRepo,
		studentRepo:  studentRepo,
	}
}

// ---------- Audience methods ----------
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

// ---------- Group methods ----------
func (s *ScheduleService) CreateGroup(ctx context.Context, group *models.Group) error {
	return s.groupRepo.Create(ctx, group)
}

func (s *ScheduleService) GetGroup(ctx context.Context, id int) (*models.Group, error) {
	return s.groupRepo.GetByID(ctx, id)
}

func (s *ScheduleService) GetAllGroups(ctx context.Context) ([]models.Group, error) {
	return s.groupRepo.GetAll(ctx)
}

func (s *ScheduleService) UpdateGroup(ctx context.Context, group *models.Group) error {
	return s.groupRepo.Update(ctx, group)
}

func (s *ScheduleService) DeleteGroup(ctx context.Context, id int) error {
	return s.groupRepo.Delete(ctx, id)
}

// ---------- Teacher methods ----------
func (s *ScheduleService) CreateTeacher(ctx context.Context, req *dto.CreateUserRequest) (*models.Teacher, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	user := models.User{
		Name:         req.Username,
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}
	if err := s.teacherRepo.CreateTeacher(ctx, &user); err != nil {
		return nil, err
	}
	return &models.Teacher{User: user}, nil
}

func (s *ScheduleService) GetTeacher(ctx context.Context, userID int) (*models.Teacher, error) {
	return s.teacherRepo.GetTeacherByUserID(ctx, userID)
}

func (s *ScheduleService) GetAllTeachers(ctx context.Context) ([]models.Teacher, error) {
	return s.teacherRepo.GetAllTeachers(ctx)
}

func (s *ScheduleService) DeleteTeacher(ctx context.Context, userID int) error {
	return s.teacherRepo.DeleteTeacher(ctx, userID)
}

// ---------- Student methods ----------
func (s *ScheduleService) CreateStudent(ctx context.Context, req *dto.CreateStudentRequest) (*models.Student, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	// Ensure group exists
	group, err := s.groupRepo.GetByID(ctx, req.GroupID)
	if err != nil {
		return nil, fmt.Errorf("get group: %w", err)
	}
	if group == nil {
		return nil, fmt.Errorf("group with id %d not found", req.GroupID)
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	user := models.User{
		Name:         req.Username,
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}
	if err := s.studentRepo.CreateStudent(ctx, &user, req.GroupID); err != nil {
		return nil, err
	}
	return &models.Student{User: user, Group: group}, nil
}

func (s *ScheduleService) GetStudent(ctx context.Context, userID int) (*models.Student, error) {
	return s.studentRepo.GetStudentByUserID(ctx, userID)
}

func (s *ScheduleService) GetAllStudents(ctx context.Context) ([]models.Student, error) {
	return s.studentRepo.GetAllStudents(ctx)
}

func (s *ScheduleService) UpdateStudentGroup(ctx context.Context, userID int, groupID int) error {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return fmt.Errorf("group with id %d not found", groupID)
	}
	return s.studentRepo.UpdateStudentGroup(ctx, userID, groupID)
}

func (s *ScheduleService) DeleteStudent(ctx context.Context, userID int) error {
	return s.studentRepo.DeleteStudent(ctx, userID)
}
