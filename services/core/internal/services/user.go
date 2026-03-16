package services

import (
	"context"
	"core/internal/dto"
	"core/internal/models"
	"core/internal/repository"
	"errors"
)

type UserService struct {
	userRepo    *repository.UserRepo
	sessionRepo *repository.SessionRepo
	jwtService  *JWTService
}

func NewUserService(userRepo *repository.UserRepo, sessionRepo *repository.SessionRepo, jwtService *JWTService) *UserService {
	return &UserService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtService:  jwtService,
	}
}

func (s *UserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if !checkPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid password")
	}

	var roles []string
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}
	accessToken, err := s.jwtService.GenerateToken(user.ID, roles)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	session := models.Session{
		User: user,
		RefreshToken: refreshToken,
		UserAgent: req.UserAgent,
		ClientIP: req.ClientIP,
	}
	err = s.sessionRepo.Create(ctx, &session)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}, nil
}
