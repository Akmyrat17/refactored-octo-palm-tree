package application

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/shared/app_errors"
	"github.com/boilerplate/pkg/logger"
	"github.com/boilerplate/pkg/query"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id domain.UserID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindAll(ctx context.Context, limit, offset int, filters []query.Filter, sorts []query.SortField) ([]*domain.User, int64, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id domain.UserID) error
}

type UserService struct {
	repository UserRepository
	logger     logger.Logger
}

func NewUserService(repository UserRepository, logger logger.Logger) *UserService {
	return &UserService{
		repository: repository,
		logger:     logger,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User, password string) error {
	// Check if email already exists
	existing, _ := s.repository.FindByEmail(ctx, user.Email)
	if existing != nil {
		return app_errors.Conflict("email is already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return app_errors.InternalError("failed to hash password").WithCause(err)
	}

	user.PasswordHash = string(hash)
	return s.repository.Save(ctx, user)
}

func (s *UserService) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *UserService) FindAll(ctx context.Context, limit, offset int, filters []query.Filter, sorts []query.SortField) ([]*domain.User, int64, error) {
	return s.repository.FindAll(ctx, limit, offset, filters, sorts)
}

func (s *UserService) UpdateUser(ctx context.Context, user *domain.User) error {
	return s.repository.Update(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id domain.UserID) error {
	return s.repository.Delete(ctx, id)
}
