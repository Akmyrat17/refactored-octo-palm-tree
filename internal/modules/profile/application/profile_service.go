package application

import (
	"context"

	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/profile/infra/persistance/dao"
	"github.com/boilerplate/pkg/query"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *dao.ProfileDAO) error
	FindByID(ctx context.Context, id domain.ProfileID) (*domain.Profile, error)
	FindAll(ctx context.Context, limit, offset int, filters []query.Filter, sorts []query.SortField) ([]*domain.Profile, int64, error)
	Update(ctx context.Context, profile *dao.ProfileDAO) error
	Delete(ctx context.Context, id domain.ProfileID) error
}
type ProfileService struct {
	profileRepo ProfileRepository
}

func NewProfileService(profileRepo ProfileRepository) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
	}
}

func (s *ProfileService) Create(ctx context.Context, profile *domain.Profile) error {
	return s.profileRepo.Create(ctx, dao.FromDomain(profile))
}

func (s *ProfileService) FindByID(ctx context.Context, id domain.ProfileID) (*domain.Profile, error) {
	return s.profileRepo.FindByID(ctx, id)
}

func (s *ProfileService) FindAll(ctx context.Context, limit, offset int, filters []query.Filter, sorts []query.SortField) ([]*domain.Profile, int64, error) {
	return s.profileRepo.FindAll(ctx, limit, offset, filters, sorts)
}

func (s *ProfileService) Update(ctx context.Context, profile *domain.Profile) error {
	return s.profileRepo.Update(ctx, dao.FromDomain(profile))
}

func (s *ProfileService) Delete(ctx context.Context, id domain.ProfileID) error {
	return s.profileRepo.Delete(ctx, id)
}
