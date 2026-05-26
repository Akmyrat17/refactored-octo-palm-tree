package domain

import (
	"time"

	"github.com/boilerplate/internal/shared/app_errors"
	"github.com/boilerplate/internal/shared/enums"
	"github.com/google/uuid"
)

type UserID uuid.UUID
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)

type User struct {
	ID           UserID
	Name         string
	Email        string
	Phone        string
	Role         enums.Role
	PasswordHash string
	TokenVersion int
	Status       UserStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(name, email, phone string) *User {
	return &User{
		ID:           UserID(uuid.New()),
		Name:         name,
		Email:        email,
		Phone:        phone,
		Role:         enums.RoleUser,
		TokenVersion: 1,
		Status:       UserStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (u UserID) String() string { return uuid.UUID(u).String() }

func ParseUserID(s string) (UserID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return UserID{}, app_errors.ValidationError("invalid user ID format")
	}
	return UserID(id), nil
}

func (u *User) Deactivate() error {
	if u.Status == UserStatusBanned {
		return nil // Cannot deactivate banned users
	}
	u.Status = UserStatusInactive
	return nil
}

func (u *User) Activate() error {
	if u.Status == UserStatusBanned {
		return nil // Cannot activate banned users
	}
	u.Status = UserStatusActive
	return nil
}

func (u *User) IncrementTokenVersion() {
	u.TokenVersion++
	u.UpdatedAt = time.Now()
}
