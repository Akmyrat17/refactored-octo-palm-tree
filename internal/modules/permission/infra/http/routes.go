package http

import (
	"github.com/boilerplate/internal/modules/permission/application"
	"github.com/boilerplate/internal/modules/permission/infra/persistence"
	"github.com/boilerplate/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group, db *pgxpool.Pool, log logger.Logger) {
	permRepo := persistence.NewPermissionRepoImpl(db)
	groupPermRepo := persistence.NewGroupPermissionRepoImpl(db)
	service := application.NewPermissionService(permRepo, groupPermRepo, log)
	handler := NewPermissionHandler(service)
	g := e.Group("/permissions")
	g.GET("", handler.ListPermissions)
}
