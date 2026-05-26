package http

import (
	"github.com/boilerplate/internal/modules/upload/application"
	"github.com/boilerplate/internal/modules/upload/domain"
	"github.com/boilerplate/internal/modules/upload/infra/http/dto"
	"github.com/boilerplate/internal/shared/response"
	"github.com/labstack/echo/v4"
)

type UploadHandler struct {
	uploadService *application.UploadService
}

func NewUploadHandler(uploadService *application.UploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

func (h *UploadHandler) UploadImage(c echo.Context) error {
	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	// Get upload type from form
	uploadTypeStr := c.FormValue("type")
	if uploadTypeStr == "" {
		uploadTypeStr = string(domain.ProfileImage) // Default to profile
	}

	uploadType, err := domain.NewUploadType(uploadTypeStr)
	if err != nil {
		return err
	}

	// Upload image
	path, err := h.uploadService.UploadImage(c.Request().Context(), file, uploadType)
	if err != nil {
		return err
	}

	return response.Created(c, dto.NewUploadImageRes(path))
}
