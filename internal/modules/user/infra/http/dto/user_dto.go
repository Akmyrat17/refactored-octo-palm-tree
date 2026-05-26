package dto

import (
	"time"

	"github.com/boilerplate/internal/domain"
)

type CreateUserReq struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

func (r *CreateUserReq) ToDomain() *domain.User {
	return domain.NewUser(r.Name, r.Email, r.Phone)
}

type UpdateUserReq struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
	Phone *string `json:"phone,omitempty"`
}

func (r *UpdateUserReq) ToDomain(existing *domain.User) {
	if r.Name != nil {
		existing.Name = *r.Name
	}
	if r.Email != nil {
		existing.Email = *r.Email
	}
	if r.Phone != nil {
		existing.Phone = *r.Phone
	}
}

type UserRes struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func UserResFromDomain(u *domain.User) UserRes {
	return UserRes{
		ID:        u.ID.String(),
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role.String(),
		Status:    string(u.Status),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func UserListResFromDomain(users []*domain.User) []UserRes {
	if len(users) == 0 {
		return make([]UserRes, 0)
	}
	res := make([]UserRes, len(users))
	for i := range users {
		res[i] = UserResFromDomain(users[i])
	}
	return res
}
