package service

import (
	"auth_service/internal/dto"
	"auth_service/internal/models"
	"auth_service/internal/repository"
	"fmt"
	"log"
	"time"
)

type RegistrationCodeService struct {
	codeRepo *repository.RegCodeRepository
	userRepo repository.UserRepository
}

func NewRegistrationCodeService(codeRepo *repository.RegCodeRepository, userRepo repository.UserRepository) *RegistrationCodeService {
	return &RegistrationCodeService{
		codeRepo: codeRepo,
		userRepo: userRepo,
	}
}

func (s *RegistrationCodeService) CreateCode(req *dto.CreateRegistrationCodeRequest) (*models.RegistrationCode, error) {
	creator, err := s.userRepo.Get(req.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("creator not found")
	}

	var allowed bool = false
	for _, role := range creator.Roles {
		if role.Name == models.AdminRole || role.Name == models.ManagerRole {
			allowed = true
			break
		}
	}

	if !allowed {
		return nil, fmt.Errorf("not allowed")
	}

	if req.RoleName == models.AdminRole {
		return nil, fmt.Errorf("not allowed")
	}

	return s.codeRepo.CreateCode(req)
}

func (s *RegistrationCodeService) ValidateCode(req *dto.ValidateCodeRequest) *dto.ValidateCodeResponse {
	regCode, err := s.codeRepo.GetCode(req.Code)
	if err != nil {
		return &dto.ValidateCodeResponse{
			Valid: false,
		}
	}

	if regCode.ExpiresAt.Before(time.Now()) {
		return &dto.ValidateCodeResponse{
			Valid: false,
		}
	}

	if regCode.CurrentUses >= regCode.MaxUses {
		return &dto.ValidateCodeResponse{
			Valid: false,
		}
	}

	return &dto.ValidateCodeResponse{
		Valid: true,
		Role: dto.RoleResponse{
			ID:          regCode.Role.ID,
			Name:        regCode.Role.Name,
			Description: regCode.Role.Description,
		},
		GroupName:     regCode.GroupName,
		RemainingUses: regCode.MaxUses - regCode.CurrentUses,
	}
}

func (s *RegistrationCodeService) RegisterWithCode(req *dto.RegisterUserRequest) (*models.User, error) {
	validation := s.ValidateCode(&dto.ValidateCodeRequest{Code: req.Code})
	if !validation.Valid {
		return nil, fmt.Errorf("code not valid")
	}

	err := s.codeRepo.UseCode(req.Code)
	if err != nil {
		return nil, err
	}
	
	user, _ := s.userRepo.Get(req.TelegramID)
	if user == nil {
		user, err = s.userRepo.Create(&models.User{
			TelegramID: req.TelegramID,
			Username: req.Username,
			FullName: req.FullName,
		})
		if err != nil {
			return nil, err
		}
	}

	roles := []string{}
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}
	roles = append(roles, validation.Role.Name)

	err = s.userRepo.UpdateUserRoles(user.TelegramID, roles)		
	if err != nil {
		log.Printf("Failed to update user roles: %v", err)
		return nil, err
	}

	// TODO: добавить изменение группы студента
	
	user, err = s.userRepo.Get(user.TelegramID)
	if err != nil {
		log.Printf("failed get updated user: %v", err)
	}


	return user, nil
}
