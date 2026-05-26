package http

import (
	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/user/application"
	"github.com/boilerplate/internal/modules/user/infra/http/dto"
	"github.com/boilerplate/internal/shared/app_errors"
	"github.com/boilerplate/internal/shared/response"
	"github.com/boilerplate/pkg/query"
	reqctx "github.com/boilerplate/pkg/req_ctx"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service *application.UserService
}

func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	p := reqctx.ParsePagination(c)
	filters := query.ParseFilters(c.QueryParams())
	sorts := query.ParseSort(c.QueryParam("sort"))
	users, total, err := h.service.FindAll(c.Request().Context(), p.PerPage, p.Offset(), filters, sorts)
	if err != nil {
		return err
	}
	return response.Paginated(c, dto.UserListResFromDomain(users), p.Page, p.PerPage, total)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req dto.CreateUserReq
	if err := reqctx.BindAndValidate(c, &req); err != nil {
		return err
	}
	user := domain.NewUser(req.Name, req.Email, req.Phone)
	if err := h.service.CreateUser(c.Request().Context(), user, req.Password); err != nil {
		return err
	}
	return response.Created(c, dto.UserResFromDomain(user))
}

func (h *UserHandler) GetUser(c echo.Context) error {
	userID, err := domain.ParseUserID(c.Param("id"))
	if err != nil {
		return app_errors.InvalidInput()
	}
	user, err := h.service.FindByID(c.Request().Context(), userID)
	if err != nil {
		return err
	}
	return response.OK(c, dto.UserResFromDomain(user))
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	userID, err := domain.ParseUserID(c.Param("id"))
	if err != nil {
		return app_errors.InvalidInput()
	}

	user, err := h.service.FindByID(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	var req dto.UpdateUserReq
	if err := reqctx.BindAndValidate(c, &req); err != nil {
		return err
	}

	req.ToDomain(user)
	if err := h.service.UpdateUser(c.Request().Context(), user); err != nil {
		return err
	}
	return response.OK(c, dto.UserResFromDomain(user))
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	userID, err := domain.ParseUserID(c.Param("id"))
	if err != nil {
		return app_errors.InvalidInput()
	}
	if err := h.service.DeleteUser(c.Request().Context(), userID); err != nil {
		return err
	}
	return response.OK(c, nil)
}
