package dao

import (
	"time"

	"github.com/boilerplate/internal/domain"
	"github.com/google/uuid"
)

type ProfileDAO struct {
	ID          uuid.UUID
	FullName    string
	AvatarURL   string
	Headline    string
	Bio         string
	ResumeURL   string
	SocialLinks map[string]string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (d *ProfileDAO) ToDomain() *domain.Profile {
	return &domain.Profile{
		ID:          domain.ProfileID(d.ID),
		FullName:    d.FullName,
		AvatarURL:   d.AvatarURL,
		Headline:    d.Headline,
		Bio:         d.Bio,
		ResumeURL:   d.ResumeURL,
		SocialLinks: d.SocialLinks,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

func FromDomain(profile *domain.Profile) *ProfileDAO {
	return &ProfileDAO{
		ID:          uuid.UUID(profile.ID),
		FullName:    profile.FullName,
		AvatarURL:   profile.AvatarURL,
		Headline:    profile.Headline,
		Bio:         profile.Bio,
		ResumeURL:   profile.ResumeURL,
		SocialLinks: profile.SocialLinks,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}
}
