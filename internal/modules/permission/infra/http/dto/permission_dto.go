package dto

import (
	"time"

	"github.com/boilerplate/internal/domain"
)

type PermissionRes struct {
	ID          int       `json:"id"`
	Module      string    `json:"module"`
	Action      string    `json:"action"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func PermissionResFromDomain(p *domain.Permission) PermissionRes {
	return PermissionRes{
		ID:          p.ID,
		Module:      p.Module,
		Action:      p.Action,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
	}
}

func PermissionListResFromDomain(permissions []*domain.Permission) []PermissionRes {
	if len(permissions) == 0 {
		return make([]PermissionRes, 0)
	}
	res := make([]PermissionRes, len(permissions))
	for i := range permissions {
		res[i] = PermissionResFromDomain(permissions[i])
	}
	return res
}

type GroupPermissionRes struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	PermissionIDs []int     `json:"permission_ids"`
	CreatedAt     time.Time `json:"created_at"`
}

func GroupPermissionResFromDomain(gp *domain.GroupPermission) GroupPermissionRes {
	return GroupPermissionRes{
		ID:            gp.ID,
		Name:          gp.Name,
		PermissionIDs: gp.PermissionIDs,
		CreatedAt:     gp.CreatedAt,
	}
}
