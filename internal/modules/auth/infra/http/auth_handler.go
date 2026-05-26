package http

import (
	"github.com/boilerplate/internal/modules/auth/application"
	"github.com/boilerplate/internal/modules/auth/infra/http/dto"
	"github.com/boilerplate/internal/shared/app_errors"
	sharedmiddleware "github.com/boilerplate/internal/shared/middleware"
	"github.com/boilerplate/internal/shared/response"
	req_ctx "github.com/boilerplate/pkg/req_ctx"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service *application.AuthService
}

func NewAuthHandler(service *application.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginReq
	if err := req_ctx.BindAndValidate(c, &req); err != nil {
		return err
	}
	pair, user, err := h.service.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return err
	}

	return response.OK(c, dto.AuthResFromPairAndUser(pair, user))
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	var req dto.RefreshReq
	if err := req_ctx.BindAndValidate(c, &req); err != nil {
		return err
	}

	pair, user, err := h.service.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return err
	}

	return response.OK(c, dto.AuthResFromPairAndUser(pair, user))
}

func (h *AuthHandler) Logout(c echo.Context) error {
	user, ok := sharedmiddleware.CurrentUser(c)
	if !ok {
		return app_errors.Unauthorized("invalid session")
	}

	if err := h.service.Logout(c.Request().Context(), user.ID); err != nil {
		return err
	}

	return response.OK(c, nil)
}
