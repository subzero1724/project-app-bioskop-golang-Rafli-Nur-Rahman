package service

import (
	"context"
	"fmt"
	"time"

	"cinema-booking-system/internal/config"
	"cinema-booking-system/internal/dto"
	"cinema-booking-system/internal/models"
	"cinema-booking-system/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication-related business logic
type AuthService struct {
	userRepo *repository.UserRepository
	config   *config.Config
	logger   *zap.Logger
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config, logger *zap.Logger) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
		logger:   logger,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*models.User, error) {
	// Check if username already exists
	existingUser, _ := s.userRepo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Check if email already exists
	existingUser, _ = s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return nil, fmt.Errorf("failed to hash password")
	}

	// Create user
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, fmt.Errorf("failed to create user")
	}

	s.logger.Info("User registered successfully", zap.String("username", user.Username))
	return user, nil
}

// Login authenticates a user and returns a token
func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		s.logger.Warn("Login attempt with invalid username", zap.String("username", req.Username))
		return nil, fmt.Errorf("invalid username or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.logger.Warn("Login attempt with invalid password", zap.String("username", req.Username))
		return nil, fmt.Errorf("invalid username or password")
	}

	// Generate JWT token
	expiresAt := time.Now().Add(s.config.GetJWTExpiration())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		s.logger.Error("Failed to sign token", zap.Error(err))
		return nil, fmt.Errorf("failed to generate token")
	}

	// Store token in database
	tokenModel := &models.Token{
		UserID:    user.ID,
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}

	if err := s.userRepo.CreateToken(ctx, tokenModel); err != nil {
		s.logger.Error("Failed to store token", zap.Error(err))
		return nil, fmt.Errorf("failed to store token")
	}

	s.logger.Info("User logged in successfully", zap.String("username", user.Username))

	return &dto.LoginResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt.Format(time.RFC3339),
		User: dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
		},
	}, nil
}

// Logout invalidates a user's token
func (s *AuthService) Logout(ctx context.Context, tokenString string) error {
	if err := s.userRepo.DeleteToken(ctx, tokenString); err != nil {
		s.logger.Error("Failed to delete token", zap.Error(err))
		return fmt.Errorf("failed to logout")
	}

	s.logger.Info("User logged out successfully")
	return nil
}

// ValidateToken verifies and validates a JWT token
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*models.User, error) {
	// Parse and validate JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Check if token exists in database
	tokenModel, err := s.userRepo.GetTokenByValue(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("token not found or expired")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, tokenModel.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}
