package reqctx

import (
	"fmt"

	"github.com/boilerplate/internal/shared/app_errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type PaginationReq struct {
	Page    int `query:"page"`
	PerPage int `query:"per_page"`
}

func ParsePagination(c echo.Context) PaginationReq {
	var p PaginationReq
	_ = c.Bind(&p)
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PerPage <= 0 {
		p.PerPage = 20
	}
	if p.PerPage > 100 {
		p.PerPage = 100
	}
	return p
}

func (p PaginationReq) Offset() int {
	return (p.Page - 1) * p.PerPage
}

func BindAndValidate(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return app_errors.InvalidInput().WithCause(err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		message := fmt.Sprintf("validation failed: %s", validationErrors[0].Field())
		return app_errors.ValidationError(message)
	}

	return nil
}
