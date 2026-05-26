package domain

import "fmt"

type UploadType string

const (
	ProfileImage UploadType = "profile"
)

var validUploadTypes = map[UploadType]bool{
	ProfileImage: true,
}

func (t UploadType) IsValid() bool {
	return validUploadTypes[t]
}

func (t UploadType) String() string {
	return string(t)
}

func NewUploadType(s string) (UploadType, error) {
	ut := UploadType(s)
	if !ut.IsValid() {
		return "", fmt.Errorf("invalid upload type: %s", s)
	}
	return ut, nil
}

var allowedImageExtensions = map[string]bool{
	"jpg":  true,
	"jpeg": true,
	"png":  true,
	"webp": true,
	"gif":  true,
}

func IsValidImageExtension(ext string) bool {
	return allowedImageExtensions[ext]
}

var mimeTypeToExt = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/webp": "webp",
	"image/gif":  "gif",
}

func ExtFromMimeType(mimeType string) string {
	if ext, ok := mimeTypeToExt[mimeType]; ok {
		return ext
	}
	return ""
}
