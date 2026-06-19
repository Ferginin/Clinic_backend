package service

import (
	"Clinic_backend/config"
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/repository"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, req *entity.UserRegisterRequest) (*entity.AuthResponse, error)
	Login(ctx context.Context, req *entity.UserLoginRequest) (*entity.AuthResponse, error)
	Refresh(ctx context.Context, req *entity.RefreshTokenRequest) (*entity.AuthResponse, error)
	Logout(ctx context.Context, req *entity.RefreshTokenRequest) error
}

type AuthService struct {
	cfg              *config.Config
	userRepo         repository.UserRepositoryInterface
	refreshTokenRepo repository.RefreshTokenRepositoryInterface
}

func NewAuthService(cfg *config.Config, userRepo repository.UserRepositoryInterface, refreshTokenRepo repository.RefreshTokenRepositoryInterface) AuthServiceInterface {
	return &AuthService{
		cfg:              cfg,
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, req *entity.UserRegisterRequest) (*entity.AuthResponse, error) {
	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return s.buildAuthResponse(ctx, createdUser)
}

func (s *AuthService) Login(ctx context.Context, req *entity.UserLoginRequest) (*entity.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user.Blocked {
		return nil, errors.New("user is blocked")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.buildAuthResponse(ctx, user)
}

func (s *AuthService) Refresh(ctx context.Context, req *entity.RefreshTokenRequest) (*entity.AuthResponse, error) {
	token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.Env.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	userID := int(userIDFloat)

	h := sha256.Sum256([]byte(req.RefreshToken))
	tokenHash := hex.EncodeToString(h[:])

	stored, err := s.refreshTokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, errors.New("refresh token not found")
	}
	if stored.RevokedAt != nil {
		return nil, errors.New("refresh token has been revoked")
	}
	if stored.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token has expired")
	}

	if err := s.refreshTokenRepo.Revoke(ctx, tokenHash); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Blocked {
		return nil, errors.New("user is blocked")
	}

	return s.buildAuthResponse(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, req *entity.RefreshTokenRequest) error {
	h := sha256.Sum256([]byte(req.RefreshToken))
	tokenHash := hex.EncodeToString(h[:])
	return s.refreshTokenRepo.Revoke(ctx, tokenHash)
}

func (s *AuthService) buildAuthResponse(ctx context.Context, user *entity.User) (*entity.AuthResponse, error) {
	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(time.Hour * time.Duration(s.cfg.Env.JWTRefreshExpireHours))
	if _, err := s.refreshTokenRepo.Create(ctx, user.ID, refreshToken, expiresAt); err != nil {
		return nil, err
	}

	return &entity.AuthResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
	}, nil
}

func (s *AuthService) generateTokens(user *entity.User) (string, string, error) {
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
