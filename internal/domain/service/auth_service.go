package service

import (
	"context"

	"github.com/bpka/tps-pkb/internal/domain/dto"
	"github.com/bpka/tps-pkb/internal/domain/port"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo port.UserRepository
	token    port.TokenProvider
}

func NewAuthService(userRepo port.UserRepository, token port.TokenProvider) port.AuthService {
	return &authService{userRepo: userRepo, token: token}
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.TokenResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, dto.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, dto.ErrInvalidCredentials
	}

	token, expiresIn, err := s.token.Generate(user)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{Token: token, ExpiresIn: expiresIn}, nil
}
