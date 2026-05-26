package dto

import "path/filepath"

type UploadImageRes struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

func NewUploadImageRes(path string) UploadImageRes {
	// Normalize path to forward slashes for URL
	normalizedPath := filepath.ToSlash(path)
	return UploadImageRes{
		Path: normalizedPath,
		URL:  "/uploads/" + normalizedPath,
	}
}
