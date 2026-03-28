package service

import (
	"context"
	"crm/internal/middleware"
	"crm/internal/model"
	repo "crm/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*retUser, error)
	Login(ctx context.Context, email, password string) (*retUser, error)
}

type authService struct {
	repo   repo.AuthRepo
	secret string
}

type retUser struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

func NewAuthService(repo repo.AuthRepo, secret string) *authService {
	return &authService{
		repo:   repo,
		secret: secret,
	}
}

func HashPassword(password string) (string, error) {
	bPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bPass), err
}

func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (s *authService) Register(ctx context.Context, email, password string) (*retUser, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	userID, err := s.repo.CreateUser(ctx, &model.User{
		Email:        email,
		PasswordHash: hash,
		Role:         "user",
	})
	if err != nil {
		return nil, err
	}
	token, err := middleware.GenerateToken(userID, "user", s.secret)
	if err != nil {
		return nil, err
	}
	return &retUser{
		Token: token,
		Role:  "user",
	}, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*retUser, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("Неверные данные")
	}
	err = CheckPassword(password, user.PasswordHash)
	if err != nil {
		return nil, errors.New("Неверные данные")
	}
	token, err := middleware.GenerateToken(user.ID, user.Role, s.secret)
	if err != nil {
		return nil, err
	}
	return &retUser{
		Token: token,
		Role:  user.Role,
	}, nil
}
