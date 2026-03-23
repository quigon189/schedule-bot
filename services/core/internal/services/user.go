package services

import (
	"context"
	"core/internal/dto"
	"core/internal/models"
	"core/internal/repository"
	"errors"
	"fmt"
)

type UserService struct {
	userRepo    *repository.UserRepo
	sessionRepo *repository.SessionRepo
	jwtService  *JWTService
	perPage     int
}

func NewUserService(userRepo *repository.UserRepo, sessionRepo *repository.SessionRepo, jwtService *JWTService) *UserService {
	return &UserService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtService:  jwtService,
		perPage:     20,
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

	refreshToken, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	session := models.Session{
		User:         user,
		RefreshToken: refreshToken,
		UserAgent:    req.UserAgent,
		ClientIP:     req.ClientIP,
	}
	err = s.sessionRepo.Create(ctx, &session)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtService.GenerateToken(user, session.ID.String())
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		SessionID:    session.ID.String(),
	}, nil
}

func (s *UserService) ValidateToken(accessToken string) (*models.User, error) {
	claims, err := s.jwtService.ValidateToken(accessToken)
	if err != nil {
		return nil, err
	}

	user := claims.User

	return &user, nil
}

func (s *UserService) ValidateSession(ctx context.Context, accessToken string) (*models.Session, error) {
	claims, err := s.jwtService.ValidateToken(accessToken)
	if err != nil {
		return nil, err
	}

	return s.sessionRepo.GetByID(ctx, claims.SessionID)
}

func (s *UserService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.LoginResponse, error) {
	session, err := s.sessionRepo.GetByID(ctx, req.SessionID)
	if err != nil {
		return nil, err
	}

	if session.RefreshToken != req.RefreshToken {
		invalidErr := errors.New("invalid refresh token")
		if err := s.sessionRepo.Delete(ctx, session.ID); err != nil {
			return nil, errors.Join(invalidErr, err)
		} else {
			return nil, invalidErr
		}
	}

	user, err := s.userRepo.GetByID(ctx, session.User.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := s.sessionRepo.UpdateRefreshToken(ctx, session.ID.String(), refreshToken); err != nil {
		return nil, err
	}

	accessToken, err := s.jwtService.GenerateToken(user, session.ID.String())
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *UserService) Logout(ctx context.Context, accessToken string) error {
	session, err := s.ValidateSession(ctx, accessToken)
	if err != nil {
		return err
	}

	return s.sessionRepo.Delete(ctx, session.ID)
}

func (s *UserService) LogoutSession(ctx context.Context, accessToken string, session *models.Session) error {
	userSession, err := s.ValidateSession(ctx, accessToken)
	if err != nil {
		return err
	}

	if userSession.User.ID != session.User.ID {
		return errors.New("access denied")
	}

	return s.sessionRepo.Delete(ctx, session.ID)
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*models.User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	password_hash, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := models.User{
		Name:         req.Username,
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: password_hash,
	}
	if err := s.userRepo.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetPaginatedUsers(ctx context.Context, req *dto.PagiantedUserRequest) (*dto.PaginatedUsers, error) {
	page := req.Page
	perPage := req.PerPage
	sortBy := req.SortBy
	sortOrder := req.SortOrder

	if page < 1 {
		page = 1
	}

	if perPage < s.perPage {
		perPage = s.perPage
	}

	if sortOrder == "" {
		sortOrder = "ASC"
	}

	if sortBy == "" {
		sortBy = "id"
	}

	return s.userRepo.GetUsersPaginated(ctx, page, perPage, sortBy, sortOrder)
}

func (s *UserService) GetUser(ctx context.Context, userID int) (*models.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *UserService) UpdatePassword(ctx context.Context, user *models.User, newPassword, oldPassword string) error {
	if !checkPasswordHash(oldPassword, user.PasswordHash) {
		return errors.New("failed to match old password")
	}

	newPasswordHash, err := hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	return s.userRepo.UpdatePasswordHash(ctx, user.ID, newPasswordHash)
}

func (s *UserService) UpdatePasswordAdmin(ctx context.Context, userID int, newPassword string) error {
	newPasswordHash, err := hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	return s.userRepo.UpdatePasswordHash(ctx, userID, newPasswordHash)

}
