package service

import (
	"Clinic_backend/config"
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	cfg      *config.Config
	userRepo *repository.UserRepository
}

func NewAuthService(cfg *config.Config, db *pgxpool.Pool) *AuthService {
	return &AuthService{
		cfg:      cfg,
		userRepo: repository.NewUserRepository(db),
	}
}

func (s *AuthService) Register(ctx context.Context, req *entity.UserRegisterRequest) (*entity.AuthResponse, error) {
	// Проверяем существование пользователя
	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаем пользователя
	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Генерируем токены
	token, refreshToken, err := s.generateTokens(createdUser)
	if err != nil {
		return nil, err
	}

	return &entity.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         createdUser.ToResponse(),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *entity.UserLoginRequest) (*entity.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user.Blocked {
		return nil, errors.New("user is blocked")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Генерируем токены
	token, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, err
	}

	return &entity.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
	}, nil
}

func (s *AuthService) generateTokens(user *entity.User) (string, string, error) {
	// Access token
	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.RoleName,
		"exp":     time.Now().Add(time.Hour * time.Duration(s.cfg.Env.JWTExpireHours)).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	token, err := accessToken.SignedString([]byte(s.cfg.Env.JWTSecret))
	if err != nil {
		return "", "", err
	}

	// Refresh token
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * time.Duration(s.cfg.Env.JWTRefreshExpireHours)).Unix(),
	}
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshTokenObj.SignedString([]byte(s.cfg.Env.JWTSecret))
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}
