package server

import (
	"fmt"

	"github.com/boilerplate/internal/config"
	authhttp "github.com/boilerplate/internal/modules/auth/infra/http"
	permissionhttp "github.com/boilerplate/internal/modules/permission/infra/http"
	profilehttp "github.com/boilerplate/internal/modules/profile/infra/http"
	uploadhttp "github.com/boilerplate/internal/modules/upload/infra/http"
	userhttp "github.com/boilerplate/internal/modules/user/infra/http"
	shared_middleware "github.com/boilerplate/internal/shared/middleware"
	"github.com/boilerplate/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(cfg *config.Config, db *pgxpool.Pool, log logger.Logger) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(shared_middleware.ErrorMiddlewareWithLogger(log))

	// Routes
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Serve uploaded files
	e.Static("/uploads", "uploads")

	e.Group("/api/v1")
	apiV1 := e.Group("/api/v1")
	authhttp.RegisterRoutes(apiV1, db, cfg, log)
	userhttp.RegisterRoutes(apiV1, db, log)
	permissionhttp.RegisterRoutes(apiV1, db, log)
	profilehttp.RegisterRoutes(apiV1, db, log)
	uploadhttp.RegisterRoutes(apiV1, log)

	return e
}

func Start(e *echo.Echo, port int) error {
	return e.Start(fmt.Sprintf(":%d", port))
}
