package application

import (
	"context"

	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/pkg/logger"
)

type PermissionRepository interface {
	FindAll(ctx context.Context) ([]*domain.Permission, error)
	SaveAll(ctx context.Context, permissions []*domain.Permission) error
}

type GroupPermissionRepository interface {
	FindByName(ctx context.Context, name string) (*domain.GroupPermission, error)
	FindAll(ctx context.Context) ([]*domain.GroupPermission, error)
	Save(ctx context.Context, gp *domain.GroupPermission) error
}

type PermissionService struct {
	permRepo      PermissionRepository
	groupPermRepo GroupPermissionRepository
	logger        logger.Logger
}

func NewPermissionService(permRepo PermissionRepository, groupPermRepo GroupPermissionRepository, logger logger.Logger) *PermissionService {
	return &PermissionService{
		permRepo:      permRepo,
		groupPermRepo: groupPermRepo,
		logger:        logger,
	}
}

func (s *PermissionService) GetAllPermissions(ctx context.Context) ([]*domain.Permission, error) {
	return s.permRepo.FindAll(ctx)
}

func (s *PermissionService) SeedPermissions(ctx context.Context, permissions []*domain.Permission) error {
	return s.permRepo.SaveAll(ctx, permissions)
}

func (s *PermissionService) GetGroupPermission(ctx context.Context, name string) (*domain.GroupPermission, error) {
	return s.groupPermRepo.FindByName(ctx, name)
}

func (s *PermissionService) SaveGroupPermission(ctx context.Context, gp *domain.GroupPermission) error {
	return s.groupPermRepo.Save(ctx, gp)
}
