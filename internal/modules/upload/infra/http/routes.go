package http

import (
	"github.com/boilerplate/internal/modules/upload/application"
	"github.com/boilerplate/pkg/logger"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group, log logger.Logger) {
	uploadService := application.NewUploadService(log)
	handler := NewUploadHandler(uploadService)
	group := e.Group("/uploads")
	group.POST("/image", handler.UploadImage)
}
