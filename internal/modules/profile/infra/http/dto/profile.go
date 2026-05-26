package dto

import (
	"time"

	"github.com/boilerplate/internal/domain"
)

type CreateProfileReq struct {
	FullName    string            `json:"full_name" validate:"required,min=2,max=100"`
	AvatarURL   string            `json:"avatar_url" validate:"omitempty"`
	Headline    string            `json:"headline" validate:"omitempty,max=255"`
	Bio         string            `json:"bio" validate:"omitempty"`
	ResumeURL   string            `json:"resume_url" validate:"omitempty,url"`
	SocialLinks map[string]string `json:"social_links" validate:"omitempty,dive,keys,required,endkeys,required"`
}

type UpdateProfileReq struct {
	FullName    *string            `json:"full_name,omitempty" validate:"omitempty,min=2,max=100"`
	AvatarURL   *string            `json:"avatar_url,omitempty" validate:"omitempty,url"`
	Headline    *string            `json:"headline,omitempty" validate:"omitempty,max=255"`
	Bio         *string            `json:"bio,omitempty" validate:"omitempty"`
	ResumeURL   *string            `json:"resume_url,omitempty" validate:"omitempty,url"`
	SocialLinks *map[string]string `json:"social_links,omitempty" validate:"omitempty,dive,keys,required,endkeys,required"`
}

type ProfileRes struct {
	ID          string            `json:"id"`
	FullName    string            `json:"full_name"`
	AvatarURL   string            `json:"avatar_url"`
	Headline    string            `json:"headline"`
	Bio         string            `json:"bio"`
	ResumeURL   string            `json:"resume_url"`
	SocialLinks map[string]string `json:"social_links"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

func ProfileResFromDomain(p *domain.Profile) ProfileRes {
	return ProfileRes{
		ID:          p.ID.String(),
		FullName:    p.FullName,
		AvatarURL:   p.AvatarURL,
		Headline:    p.Headline,
		Bio:         p.Bio,
		ResumeURL:   p.ResumeURL,
		SocialLinks: p.SocialLinks,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
	}
}

func ProfileResFromDomainList(profiles []*domain.Profile) []ProfileRes {
	res := make([]ProfileRes, len(profiles))
	for i, p := range profiles {
		res[i] = ProfileResFromDomain(p)
	}
	return res
}

func ParseProfileID(s string) (domain.ProfileID, error) {
	id, err := domain.ParseProfileID(s)
	if err != nil {
		return domain.ProfileID{}, err
	}
	return id, nil
}

func (r *CreateProfileReq) ToDomain() *domain.Profile {
	return domain.NewProfile(r.FullName, r.AvatarURL, r.Headline, r.Bio, r.ResumeURL, r.SocialLinks)
}

func (r *UpdateProfileReq) ToDomain(existing *domain.Profile) *domain.Profile {
	if r.FullName != nil {
		existing.FullName = *r.FullName
	}
	if r.AvatarURL != nil {
		existing.AvatarURL = *r.AvatarURL
	}
	if r.Headline != nil {
		existing.Headline = *r.Headline
	}
	if r.Bio != nil {
		existing.Bio = *r.Bio
	}
	if r.ResumeURL != nil {
		existing.ResumeURL = *r.ResumeURL
	}
	if r.SocialLinks != nil {
		existing.SocialLinks = *r.SocialLinks
	}
	return existing
}
