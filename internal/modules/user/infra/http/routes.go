package http

import (
	"github.com/boilerplate/internal/modules/user/application"
	"github.com/boilerplate/internal/modules/user/infra/persistence"
	"github.com/boilerplate/internal/shared/middleware"
	"github.com/boilerplate/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group, db *pgxpool.Pool, log logger.Logger) {
	repo := persistence.NewUserRepoImpl(db)
	service := application.NewUserService(repo, log)
	handler := NewUserHandler(service)

	g := e.Group("/users", middleware.Auth())
	g.GET("", handler.ListUsers)
	g.POST("", handler.CreateUser)
	g.GET("/:id", handler.GetUser)
	g.PATCH("/:id", handler.UpdateUser, middleware.RequirePermission("users.update"))
	g.DELETE("/:id", handler.DeleteUser, middleware.RequirePermission("users.delete"))
}
