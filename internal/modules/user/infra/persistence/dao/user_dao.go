package dao

import (
	"time"

	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/shared/enums"
	"github.com/google/uuid"
)

type UserDAO struct {
	ID           uuid.UUID
	Name         string
	Email        string
	Phone        string
	Role         string
	PasswordHash string
	TokenVersion int
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (d *UserDAO) ToDomain() *domain.User {
	return &domain.User{
		ID:           domain.UserID(d.ID),
		Name:         d.Name,
		Email:        d.Email,
		Phone:        d.Phone,
		Role:         enums.Role(d.Role),
		PasswordHash: d.PasswordHash,
		TokenVersion: d.TokenVersion,
		Status:       domain.UserStatus(d.Status),
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

func FromDomain(u *domain.User) UserDAO {
	return UserDAO{
		ID:           uuid.UUID(u.ID),
		Name:         u.Name,
		Email:        u.Email,
		Phone:        u.Phone,
		Role:         u.Role.String(),
		PasswordHash: u.PasswordHash,
		TokenVersion: u.TokenVersion,
		Status:       string(u.Status),
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
