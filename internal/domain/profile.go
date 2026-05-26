package domain

import (
	"time"

	"github.com/boilerplate/internal/shared/app_errors"
	"github.com/google/uuid"
)

type ProfileID uuid.UUID

func (p ProfileID) String() string { return uuid.UUID(p).String() }
func ParseProfileID(s string) (ProfileID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return ProfileID{}, app_errors.ValidationError("invalid profile ID format")
	}
	return ProfileID(id), nil
}

type Profile struct {
	ID          ProfileID
	FullName    string
	AvatarURL   string
	Headline    string
	Bio         string
	ResumeURL   string
	SocialLinks map[string]string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProfile(fullName, avatarURL, headline, bio, resumeURL string, socialLinks map[string]string) *Profile {
	return &Profile{
		ID:          ProfileID(uuid.New()),
		FullName:    fullName,
		AvatarURL:   avatarURL,
		Headline:    headline,
		Bio:         bio,
		ResumeURL:   resumeURL,
		SocialLinks: socialLinks,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
