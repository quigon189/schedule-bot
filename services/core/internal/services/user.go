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
		if err := s.sessionRepo.Delete(ctx, req.SessionID); err != nil {
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
		return nil
	}

	return s.sessionRepo.Delete(ctx, session.ID.String())
}

func (s *UserService) LogoutSession(ctx context.Context, accessToken string, session *models.Session) error {
	userSession, err := s.ValidateSession(ctx, accessToken)
	if err != nil {
		return err
	}

	if userSession.User.ID != session.User.ID {
		return errors.New("access denied")
	}

	return s.sessionRepo.Delete(ctx, session.ID.String())
}
