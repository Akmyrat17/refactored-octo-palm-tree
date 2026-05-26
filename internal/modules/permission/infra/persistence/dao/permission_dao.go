package dao

import (
	"time"

	"github.com/boilerplate/internal/domain"
)

type PermissionDAO struct {
	ID          int
	Module      string
	Action      string
	Name        string
	Description string
	CreatedAt   time.Time
}

func (d *PermissionDAO) ToDomain() *domain.Permission {
	return &domain.Permission{
		ID:          d.ID,
		Module:      d.Module,
		Action:      d.Action,
		Name:        d.Name,
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
	}
}

func FromDomain(p *domain.Permission) PermissionDAO {
	return PermissionDAO{
		ID:          p.ID,
		Module:      p.Module,
		Action:      p.Action,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
	}
}

type GroupPermissionDAO struct {
	ID            int
	Name          string
	PermissionIDs []int
	CreatedAt     time.Time
}

func (d *GroupPermissionDAO) ToDomain() *domain.GroupPermission {
	return &domain.GroupPermission{
		ID:            d.ID,
		Name:          d.Name,
		PermissionIDs: d.PermissionIDs,
		CreatedAt:     d.CreatedAt,
	}
}

func FromGroupDomain(gp *domain.GroupPermission) GroupPermissionDAO {
	return GroupPermissionDAO{
		ID:            gp.ID,
		Name:          gp.Name,
		PermissionIDs: gp.PermissionIDs,
		CreatedAt:     gp.CreatedAt,
	}
}
