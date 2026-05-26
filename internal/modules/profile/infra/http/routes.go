package http

import (
	"github.com/boilerplate/internal/modules/profile/application"
	"github.com/boilerplate/internal/modules/profile/infra/persistance"
	"github.com/boilerplate/internal/shared/middleware"
	"github.com/boilerplate/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group, db *pgxpool.Pool, log logger.Logger) {
	profileRepo := persistance.NewProfileRepoImpl(db)
	profileService := application.NewProfileService(profileRepo)
	handler := NewProfileHandler(profileService)
	group := e.Group("/profiles")
	group.Use(middleware.Auth(), middleware.IsAdmin())
	group.POST("", handler.CreateProfile)
	group.GET("/:id", handler.GetProfile)
	group.PATCH("/:id", handler.UpdateProfile)
	group.DELETE("/:id", handler.DeleteProfile)
	group.GET("", handler.ListProfiles)
}
