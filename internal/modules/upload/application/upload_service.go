package application

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/boilerplate/internal/modules/upload/domain"
	"github.com/boilerplate/internal/shared/app_errors"
	"github.com/boilerplate/pkg/logger"
)

const (
	maxFileSize = 5 * 1024 * 1024 // 5MB
	uploadsDir  = "uploads"
)

type UploadService struct {
	uploadsDir string
	log        logger.Logger
}

func NewUploadService(log logger.Logger) *UploadService {
	return &UploadService{
		uploadsDir: uploadsDir,
		log:        log,
	}
}

func (s *UploadService) UploadImage(ctx context.Context, file *multipart.FileHeader, uploadType domain.UploadType) (string, error) {
	// Validate upload type
	if !uploadType.IsValid() {
		s.log.Error("invalid upload type", "type", uploadType.String())
		return "", app_errors.InvalidInput()
	}

	// Check file size
	if file.Size > maxFileSize {
		s.log.Warn("file too large", "size", file.Size, "max", maxFileSize)
		return "", app_errors.ValidationError(fmt.Sprintf("file size exceeds %dMB limit", maxFileSize/1024/1024))
	}

	s.log.Info("starting file upload", "filename", file.Filename, "size", file.Size, "type", uploadType.String())

	// Open file
	src, err := file.Open()
	if err != nil {
		s.log.Error("failed to open file", "filename", file.Filename, "error", err)
		return "", app_errors.InvalidInput().WithCause(err)
	}
	defer src.Close()

	// Validate mime type
	mimeType, err := s.validateMimeType(src)
	if err != nil {
		s.log.Error("mime type validation failed", "filename", file.Filename, "error", err)
		return "", err
	}

	s.log.Debug("mime type detected", "mime_type", mimeType)

	// Get extension from mime type
	ext := domain.ExtFromMimeType(mimeType)
	if ext == "" {
		s.log.Error("unsupported mime type", "mime_type", mimeType)
		return "", app_errors.ValidationError("unsupported image format")
	}

	// Create upload directory if not exists
	uploadPath := filepath.Join(s.uploadsDir, uploadType.String())
	s.log.Debug("creating upload directory", "path", uploadPath)
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		s.log.Error("failed to create upload directory", "path", uploadPath, "error", err)
		return "", app_errors.InternalError("failed to create upload directory").WithCause(err)
	}

	// Generate filename with timestamp
	filename := fmt.Sprintf("%d.%s", time.Now().UnixMilli(), ext)
	filePath := filepath.Join(uploadPath, filename)

	s.log.Debug("creating destination file", "path", filePath)
	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		s.log.Error("failed to create file", "path", filePath, "error", err)
		return "", app_errors.InternalError("failed to create file").WithCause(err)
	}
	defer dst.Close()

	// Copy file
	s.log.Debug("copying file content")
	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(filePath) // Clean up on error
		s.log.Error("failed to save file", "path", filePath, "error", err)
		return "", app_errors.InternalError("failed to save file").WithCause(err)
	}

	// Return relative path
	relativePath := filepath.Join(uploadType.String(), filename)
	relativePath = filepath.ToSlash(relativePath) // Use forward slashes for consistency
	s.log.Info("file uploaded successfully", "path", relativePath)
	return relativePath, nil
}

func (s *UploadService) validateMimeType(file multipart.File) (string, error) {
	// Read first 512 bytes to detect mime type
	header := make([]byte, 512)
	n, err := file.Read(header)
	if err != nil && err != io.EOF {
		s.log.Error("failed to read file header", "error", err)
		return "", app_errors.InvalidInput().WithCause(err)
	}

	// Reset file pointer
	file.Seek(0, 0)

	// Detect mime type
	mimeType := detectMimeType(header[:n])

	// Validate it's an image
	if !strings.HasPrefix(mimeType, "image/") {
		s.log.Warn("file is not an image", "mime_type", mimeType)
		return "", app_errors.ValidationError("file must be an image")
	}

	return mimeType, nil
}

func detectMimeType(data []byte) string {
	// Simple mime type detection based on magic numbers
	if len(data) >= 4 {
		if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
			return "image/jpeg"
		}
		if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
			return "image/png"
		}
		if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 {
			return "image/gif"
		}
		if data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46 {
			return "image/webp"
		}
	}
	return "application/octet-stream"
}
