package http

import (
	"net/http"

	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/profile/application"
	"github.com/boilerplate/internal/modules/profile/infra/http/dto"
	"github.com/boilerplate/internal/shared/response"
	"github.com/boilerplate/pkg/query"
	reqctx "github.com/boilerplate/pkg/req_ctx"
	"github.com/labstack/echo/v4"
)

type ProfileHandler struct {
	profileService *application.ProfileService
}

func NewProfileHandler(profileService *application.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

func (h *ProfileHandler) CreateProfile(c echo.Context) error {
	var req dto.CreateProfileReq
	if err := reqctx.BindAndValidate(c, &req); err != nil {
		return err
	}
	profile := req.ToDomain()
	if err := h.profileService.Create(c.Request().Context(), profile); err != nil {
		return err
	}
	return response.Created(c, dto.ProfileResFromDomain(profile))
}

func (h *ProfileHandler) GetProfile(c echo.Context) error {
	profileID, err := domain.ParseProfileID(c.Param("id"))
	if err != nil {
		return err
	}
	profile, err := h.profileService.FindByID(c.Request().Context(), profileID)
	if err != nil {
		return err
	}
	return response.OK(c, dto.ProfileResFromDomain(profile))
}

func (h *ProfileHandler) UpdateProfile(c echo.Context) error {
	profileID, err := domain.ParseProfileID(c.Param("id"))
	if err != nil {
		return err
	}
	var req dto.UpdateProfileReq
	if err := reqctx.BindAndValidate(c, &req); err != nil {
		return err
	}
	existingProfile, err := h.profileService.FindByID(c.Request().Context(), profileID)
	if err != nil {
		return err
	}
	profile := req.ToDomain(existingProfile)
	profile.ID = profileID
	if err := h.profileService.Update(c.Request().Context(), profile); err != nil {
		return err
	}
	return response.OK(c, dto.ProfileResFromDomain(profile))
}

func (h *ProfileHandler) DeleteProfile(c echo.Context) error {
	profileID, err := domain.ParseProfileID(c.Param("id"))
	if err != nil {
		return err
	}
	if err := h.profileService.Delete(c.Request().Context(), profileID); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *ProfileHandler) ListProfiles(c echo.Context) error {
	var pagination reqctx.PaginationReq = reqctx.ParsePagination(c)
	filters := query.ParseFilters(c.QueryParams())
	sorts := query.ParseSort(c.QueryParam("sort"))
	profiles, total, err := h.profileService.FindAll(c.Request().Context(), pagination.PerPage, pagination.Offset(), filters, sorts)
	if err != nil {
		return err
	}
	return response.Paginated(c, dto.ProfileResFromDomainList(profiles), pagination.Page, pagination.PerPage, total)
}
