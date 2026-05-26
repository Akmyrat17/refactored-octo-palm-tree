package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/shared/app_errors"
	"github.com/boilerplate/internal/shared/enums"
	"github.com/boilerplate/internal/shared/response"
	"github.com/boilerplate/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
)

func ErrorMiddlewareWithLogger(log logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			// Handle Echo HTTP errors
			var httpErr *echo.HTTPError
			if errors.As(err, &httpErr) {
				log.Error("http error occurred",
					"code", httpErr.Code,
					"message", httpErr.Message,
					"path", c.Request().URL.Path,
					"method", c.Request().Method,
					"internal_error", fmt.Sprintf("%v", httpErr.Internal))
				return response.Error(c, app_errors.ErrCodeInternal, httpErr.Code, fmt.Sprintf("%v", httpErr.Message))
			}

			var appErr *app_errors.AppError
			if errors.As(err, &appErr) {
				// Log app errors with their cause if present
				var message string
				if appErr.Cause != nil {
					message = appErr.Cause.Error()
					log.Error("app error occurred", "code", appErr.Code, "message", appErr.Message, "cause", message, "path", c.Request().URL.Path, "method", c.Request().Method)
				} else {
					message = appErr.Message
					log.Error("app error occurred", "code", appErr.Code, "message", message, "path", c.Request().URL.Path, "method", c.Request().Method)
				}
				return response.Error(c, appErr.Code, appErr.Status, message)
			}

			// Handle pgx errors
			if errors.Is(err, pgx.ErrNoRows) {
				log.Error("resource not found", "error", err.Error(), "path", c.Request().URL.Path, "method", c.Request().Method)
				return response.Error(c, app_errors.ErrCodeNotFound, http.StatusNotFound, "resource not found")
			}

			// Handle PostgreSQL errors
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				errMsg := fmt.Sprintf("database error: %s (code: %s)", pgErr.Message, pgErr.Code)
				log.Error("postgres error occurred", "code", pgErr.Code, "message", pgErr.Message, "detail", pgErr.Detail, "path", c.Request().URL.Path, "method", c.Request().Method)
				return response.Error(c, app_errors.ErrCodeDatabase, http.StatusInternalServerError, errMsg)
			}

			// Log unknown errors with full stack
			log.Error("unexpected error occurred", "error", err.Error(), "path", c.Request().URL.Path, "method", c.Request().Method, "error_type", fmt.Sprintf("%T", err))

			// Default to internal error
			return response.Error(c, app_errors.ErrCodeInternal, http.StatusInternalServerError, err.Error())
		}
	}
}

// Deprecated: Use ErrorMiddlewareWithLogger instead
func ErrorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		var appErr *app_errors.AppError
		if errors.As(err, &appErr) {
			return response.Error(c, appErr.Code, appErr.Status, appErr.Message)
		}

		// Handle pgx errors
		if errors.Is(err, pgx.ErrNoRows) {
			return response.Error(c, app_errors.ErrCodeNotFound, http.StatusNotFound, "resource not found")
		}

		// Default to internal error
		return response.Error(c, app_errors.ErrCodeInternal, http.StatusInternalServerError, "an internal error occurred")
	}
}

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Request logging could go here
		start := time.Now()

		err := next(c)

		latency := time.Since(start).Milliseconds()
		_ = latency // Used in logging

		// Response logging could go here
		return err
	}
}

type TokenAuthenticator interface {
	ValidateToken(ctx context.Context, token string) (*domain.User, error)
}

var authProvider TokenAuthenticator

func SetAuthProvider(provider TokenAuthenticator) {
	authProvider = provider
}

func Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return app_errors.Unauthorized("authorization header missing")
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				return app_errors.Unauthorized("invalid authorization header")
			}

			token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			if token == "" {
				return app_errors.Unauthorized("missing bearer token")
			}

			if authProvider == nil {
				return app_errors.InternalError("auth provider is not configured")
			}

			user, err := authProvider.ValidateToken(c.Request().Context(), token)
			if err != nil {
				return err
			}

			c.Set("user", user)
			return next(c)
		}
	}
}

func CurrentUser(c echo.Context) (*domain.User, bool) {
	user, ok := c.Get("user").(*domain.User)
	return user, ok
}

func IsAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := CurrentUser(c)
			if !ok {
				return app_errors.Unauthorized("invalid session")
			}

			if user.Role != enums.RoleAdmin {
				return app_errors.Forbidden("admin access required")
			}
			return next(c)
		}
	}
}

func RequirePermission(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, ok := CurrentUser(c)
			if !ok {
				return app_errors.Unauthorized("invalid session")
			}

			return next(c)
		}
	}
}
