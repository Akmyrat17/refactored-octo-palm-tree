package application

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/boilerplate/internal/config"
	"github.com/boilerplate/internal/domain"
	userapp "github.com/boilerplate/internal/modules/user/application"
	"github.com/boilerplate/internal/shared/app_errors"
	"github.com/boilerplate/pkg/logger"
)

// TokenPair represents the access and refresh token pair returned to the client
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// TokenClaims represents the JWT claims structure
type TokenClaims struct {
	UserID       string `json:"user_id"`
	Role         string `json:"role"`
	TokenVersion int    `json:"token_version"`
	TokenType    string `json:"token_type"`
	jwt.RegisteredClaims
}

type AuthService struct {
	userRepo           userapp.UserRepository
	jwtSecret          string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
	logger             logger.Logger
}

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

func NewAuthService(userRepo userapp.UserRepository, cfg config.JWTConfig, logger logger.Logger) *AuthService {
	return &AuthService{
		userRepo:           userRepo,
		jwtSecret:          cfg.Secret,
		accessTokenExpiry:  cfg.AccessTokenExpiration,
		refreshTokenExpiry: cfg.RefreshTokenExpiration,
		logger:             logger,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*TokenPair, *domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if appErr, ok := err.(*app_errors.AppError); ok && appErr.Code == app_errors.ErrCodeNotFound {
			return nil, nil, app_errors.Unauthorized("invalid credentials")
		}
		s.logger.Error("Failed to find user by email", "email", email, "error", err)
		return nil, nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, nil, app_errors.Unauthorized("invalid credentials")
	}

	if user.Status != domain.UserStatusActive {
		return nil, nil, app_errors.Unauthorized("user is not active")
	}

	pair, err := s.generateTokenPair(user)
	if err != nil {
		return nil, nil, app_errors.InternalError("failed to generate tokens").WithCause(err)
	}

	s.logger.DBInfo(ctx, "user login", "user_id", user.ID.String())
	return pair, user, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, *domain.User, error) {
	user, err := s.validateToken(ctx, refreshToken, TokenTypeRefresh)
	if err != nil {
		return nil, nil, err
	}
	user.IncrementTokenVersion()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, nil, app_errors.InternalError("failed to refresh token").WithCause(err)
	}

	pair, err := s.generateTokenPair(user)
	if err != nil {
		return nil, nil, app_errors.InternalError("failed to generate tokens").WithCause(err)
	}

	return pair, user, nil
}

func (s *AuthService) Logout(ctx context.Context, userID domain.UserID) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.IncrementTokenVersion()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return app_errors.InternalError("failed to logout").WithCause(err)
	}

	s.logger.DBInfo(ctx, "user logout", "user_id", user.ID.String())
	return nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*domain.User, error) {
	user, err := s.validateToken(ctx, token, TokenTypeAccess)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) validateToken(ctx context.Context, tokenString, expectedType string) (*domain.User, error) {
	claims := &TokenClaims{}
	parser := &jwt.Parser{}
	if _, _, err := parser.ParseUnverified(tokenString, claims); err != nil {
		return nil, app_errors.Unauthorized("invalid token")
	}

	userID, err := domain.ParseUserID(claims.UserID)
	if err != nil {
		return nil, app_errors.Unauthorized("invalid token")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, app_errors.Unauthorized("invalid token")
	}

	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.signingKey(user.TokenVersion)), nil
	})
	if err != nil || !parsedToken.Valid {
		return nil, app_errors.Unauthorized("invalid token")
	}

	if claims.TokenType != expectedType {
		return nil, app_errors.Unauthorized("invalid token")
	}

	if user.Status != domain.UserStatusActive {
		return nil, app_errors.Unauthorized("user is not active")
	}

	return user, nil
}

func (s *AuthService) generateTokenPair(user *domain.User) (*TokenPair, error) {
	accessToken, err := s.buildToken(user, TokenTypeAccess, s.accessTokenExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.buildToken(user, TokenTypeRefresh, s.refreshTokenExpiry)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.accessTokenExpiry.Seconds()),
	}, nil
}

func (s *AuthService) buildToken(user *domain.User, tokenType string, expiresIn time.Duration) (string, error) {
	claims := TokenClaims{
		UserID:       user.ID.String(),
		Role:         user.Role.String(),
		TokenVersion: user.TokenVersion,
		TokenType:    tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.signingKey(user.TokenVersion)))
}

func (s *AuthService) signingKey(version int) string {
	return fmt.Sprintf("%s|%d", s.jwtSecret, version)
}
