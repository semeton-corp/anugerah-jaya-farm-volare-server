package rest

import (
	"strconv"

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

type PresenceHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IPresenceService
}

func (a *PresenceHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/presences")

	v1.Get("/current", middleware.Authentication(), a.GetCurrentStaffPresence)
	v1.Get("/", middleware.Authentication(), a.GetAllStaffPresences)
	v1.Patch("/arrival/:id", middleware.Authentication(), a.ArrivalPresence)
	v1.Patch("/departure/:id", middleware.Authentication(), a.DeparturePresence)
}

func NewPresenceHandler(log *zap.Logger, service service.IPresenceService, validator *validator.Validate) *PresenceHandler {
	return &PresenceHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (a *PresenceHandler) GetCurrentStaffPresence(c *fiber.Ctx) error {
	accountIdCtx := c.Locals("accountId").(string)
	if accountIdCtx == "" {
		a.log.Error("[GetCurrentStaffPresence] accountId not found in context")
		return errx.NotFound("accountId not found in context")
	}

	staffPresence, err := a.service.GetCurrentStaffPresence(uuid.MustParse(accountIdCtx))
	if err != nil {
		a.log.Error("[GetCurrentStaffPresence] failed to get current staff presence", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, staffPresence, "success get current staff presence")
}

func (a *PresenceHandler) GetAllStaffPresences(c *fiber.Ctx) error {
	accountIdCtx := c.Locals("accountId").(string)
	if accountIdCtx == "" {
		a.log.Error("[GetAllStaffPresences] accountId not found in context")
		return errx.NotFound("accountId not found in context")
	}

	var filter dto.GetPresenceFilter
	if err := c.QueryParser(&filter); err != nil {
		a.log.Error("[GetAllStaffPresences] failed to parsing query filter", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(filter); err != nil {
		a.log.Error("[GetAllStaffPresences] failed to validate filter", zap.Error(err))
		return err
	}

	staffPresences, err := a.service.GetAllStaffPresences(uuid.MustParse(accountIdCtx), filter)
	if err != nil {
		a.log.Error("[GetAllStaffPresences] failed to get all staff presences", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, staffPresences, "success get all staff presences")
}

func (a *PresenceHandler) ArrivalPresence(c *fiber.Ctx) error {
	accountIdCtx := c.Locals("accountId").(string)
	if accountIdCtx == "" {
		a.log.Error("[ArrivalPresence] accountId not found in context")
		return errx.NotFound("accountId not found in context")
	}

	presenceIdStr := c.Params("id")
	if presenceIdStr == "" {
		a.log.Error("[ArrivalPresence] presenceId not found in param")
		return errx.NotFound("presenceId not found in param")
	}

	presenceId, err := strconv.ParseUint(presenceIdStr, 10, 64)
	if err != nil {
		a.log.Error("[ArrivalPresence] failed to parse presenceId", zap.Error(err))
		return errx.BadRequest("presenceId not valid")
	}

	staffPresence, err := a.service.ArrivalPresence(presenceId, uuid.MustParse(accountIdCtx))
	if err != nil {
		a.log.Error("[ArrivalPresence] failed to arrival presence", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, staffPresence, "success arrival presence")
}

func (a *PresenceHandler) DeparturePresence(c *fiber.Ctx) error {
	accountIdCtx := c.Locals("accountId").(string)
	if accountIdCtx == "" {
		a.log.Error("[DeparturePresence] accountId not found in context")
		return errx.NotFound("accountId not found in context")
	}

	presenceIdStr := c.Params("id")
	if presenceIdStr == "" {
		a.log.Error("[DeparturePresence] presenceId not found in param")
		return errx.NotFound("presenceId not found in param")
	}

	presenceId, err := strconv.ParseUint(presenceIdStr, 10, 64)
	if err != nil {
		a.log.Error("[DeparturePresence] failed to parse presenceId", zap.Error(err))
		return errx.BadRequest("presenceId not valid")
	}

	staffPresence, err := a.service.DeparturePresence(presenceId, uuid.MustParse(accountIdCtx))
	if err != nil {
		a.log.Error("[DeparturePresence] failed to departure presence", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, staffPresence, "success departure presence")
}
