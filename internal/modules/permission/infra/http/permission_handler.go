package http

import (
	"github.com/boilerplate/internal/modules/permission/application"
	"github.com/boilerplate/internal/modules/permission/infra/http/dto"
	"github.com/boilerplate/internal/shared/response"
	"github.com/labstack/echo/v4"
)

type PermissionHandler struct {
	service *application.PermissionService
}

func NewPermissionHandler(service *application.PermissionService) *PermissionHandler {
	return &PermissionHandler{service: service}
}

func (h *PermissionHandler) ListPermissions(c echo.Context) error {
	permissions, err := h.service.GetAllPermissions(c.Request().Context())
	if err != nil {
		return err
	}
	return response.OK(c, dto.PermissionListResFromDomain(permissions))
}
