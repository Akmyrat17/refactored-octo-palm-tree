package http

import (
	"github.com/boilerplate/internal/config"
	"github.com/boilerplate/internal/modules/auth/application"
	userPersistence "github.com/boilerplate/internal/modules/user/infra/persistence"
	sharedmiddleware "github.com/boilerplate/internal/shared/middleware"
	"github.com/boilerplate/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group, db *pgxpool.Pool, cfg *config.Config, log logger.Logger) {
	userRepo := userPersistence.NewUserRepoImpl(db)
	authService := application.NewAuthService(userRepo, cfg.JWT, log)
	sharedmiddleware.SetAuthProvider(authService)

	handler := NewAuthHandler(authService)
	authGroup := e.Group("/auth")
	authGroup.POST("/login", handler.Login)
	authGroup.POST("/refresh", handler.Refresh)
	authGroup.POST("/logout", handler.Logout, sharedmiddleware.Auth())
}
