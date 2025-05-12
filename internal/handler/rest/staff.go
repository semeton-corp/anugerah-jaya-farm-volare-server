package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/middleware"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/service"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/response"
	"go.uber.org/zap"
)

type StaffHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IStaffService
}

func (a *StaffHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/staffs")

	v1.Put("/me", middleware.Authentication(), a.UpdateOwnProfile)
	v1.Get("/me", middleware.Authentication(), a.GetOwnProfile)
	v1.Get("", middleware.Authentication(), a.GetStaffs)
	v1.Get("/:id", middleware.Authentication(), a.GetStaffById)
	v1.Put("/:id", middleware.Authentication(), a.UpdateStaff)
	v1.Get("/overview/:id", middleware.Authentication(), a.GetOverviewStaff)
}

func NewStaffHandler(log *zap.Logger, service service.IStaffService, validator *validator.Validate) *StaffHandler {
	return &StaffHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (a *StaffHandler) GetStaffById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		a.log.Error("[GetStaffById] id param not found")
		return errx.BadRequest("id param not found")
	}

	resp, err := a.service.GetStaffById(uuid.MustParse(idParam))
	if err != nil {
		a.log.Error("[GetStaffById] failed to get staff by id")
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get staff by id")
}

func (a *StaffHandler) GetOwnProfile(c *fiber.Ctx) error {
	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("account id not found in context")
		return errx.BadRequest("account id not found in context")
	}

	resp, err := a.service.GetStaffById(uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[GetStaffById] failed to get staff by id")
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get staff by id")
}

func (a *StaffHandler) UpdateStaff(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		a.log.Error("[GetStaffById] id param not found")
		return errx.BadRequest("id param not found")
	}

	var request dto.UpdateStaffRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[UpdateAccount] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[UpdateAccount] failed to validate request", zap.Error(err))
		return err
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[UpdateAccount] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	res, err := a.service.UpdateStaff(uuid.MustParse(idParam), request, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[UpdateAccount] failed to update account", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		res,
		"success update account",
	)
}

func (a *StaffHandler) UpdateOwnProfile(c *fiber.Ctx) error {
	var request dto.UpdateStaffRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[UpdateAccount] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[UpdateAccount] failed to validate request", zap.Error(err))
		return err
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[UpdateAccount] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	res, err := a.service.UpdateStaff(uuid.MustParse(accountId), request, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[UpdateAccount] failed to update account", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		res,
		"success update account",
	)
}

func (a *StaffHandler) GetStaffs(c *fiber.Ctx) error {
	var filter dto.GetStaffFilter
	if err := c.QueryParser(&filter); err != nil {
		a.log.Error("[GetStaffs] failed to parsing query filter", zap.Error(err))
		return err
	}

	resp, err := a.service.GetStaffs(filter)
	if err != nil {
		a.log.Error("[GetStaffs] failed to get staffs")
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success to get staffs")
}

func (a *StaffHandler) GetOverviewStaff(c *fiber.Ctx) error {
	var filter dto.GetStaffOverviewFilter
	if err := c.QueryParser(&filter); err != nil {
		a.log.Error("[GetOverviewStaff] failed to parse query", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(filter); err != nil {
		a.log.Error("[GetOverviewStaff] failed to validate query", zap.Error(err))
		return err
	}

	idParam := c.Params("id")
	if idParam == "" {
		a.log.Error("[GetOverviewStaff] id param is empty")
		return errx.BadRequest("id param is empty")
	}

	resp, err := a.service.GetOverviewStaff(uuid.MustParse(idParam), filter)
	if err != nil {
		a.log.Error("[GetOverviewStaff] failed to get overview staff", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success to get overview staff")
}
