package service

import (
	"auth_service/internal/dto"
	"auth_service/internal/models"
	"auth_service/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(req *dto.CreateUserRequest) (*models.User, error) {
	existingUser, _ := s.userRepo.Get(req.TelegramID)
	if existingUser != nil {
		return existingUser, nil
	}

	user := models.User{
		TelegramID: req.TelegramID,
		Username: req.Username,
		FullName: req.FullName,
	}

	return s.userRepo.Create(&user)
}

func (s *UserService) GetUser(telegramID int64) (*models.User, error) {
	return s.userRepo.Get(telegramID)
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

func (s *UserService) UpdateUserRoles(user *models.User) error {
	var roles []string
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}
	return s.userRepo.UpdateUserRoles(user.TelegramID, roles)
}

func (s *UserService) DeleteUser(user *models.User) error {
	return s.userRepo.Delete(user.TelegramID)
}

func(s *UserService) ListUsers() ([]models.User, error) {
	return s.userRepo.List()
}
