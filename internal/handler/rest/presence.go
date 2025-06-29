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

func (h *PresenceHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/presences")

	v1.Get("/current", middleware.Authentication(), h.GetCurrentStaffPresence)
	v1.Get("/", middleware.Authentication(), h.GetAllStaffPresences)
	v1.Patch("/arrival/:id", middleware.Authentication(), h.ArrivalPresence)
	v1.Patch("/departure/:id", middleware.Authentication(), h.DeparturePresence)
}

func NewPresenceHandler(log *zap.Logger, service service.IPresenceService, validator *validator.Validate) *PresenceHandler {
	return &PresenceHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *PresenceHandler) GetCurrentStaffPresence(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[GetCurrentStaffPresence] userId not found in context")
		return errx.NotFound("userId not found in context")
	}

	staffPresence, err := h.service.GetCurrentStaffPresence(uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[GetCurrentStaffPresence] failed to get current staff presence", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, staffPresence, "success get current staff presence")
}

func (h *PresenceHandler) GetAllStaffPresences(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[GetAllStaffPresences] userId not found in context")
		return errx.NotFound("userId not found in context")
	}

	var filter dto.GetPresenceFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("[GetAllStaffPresences] failed to parsing query filter", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("[GetAllStaffPresences] failed to validate filter", zap.Error(err))
		return err
	}

	staffPresences, err := h.service.GetAllStaffPresences(uuid.MustParse(userId), filter)
	if err != nil {
		h.log.Error("[GetAllStaffPresences] failed to get all staff presences", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, staffPresences, "success get all staff presences")
}

func (h *PresenceHandler) ArrivalPresence(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[ArrivalPresence] userId not found in context")
		return errx.NotFound("userId not found in context")
	}

	presenceIdStr := c.Params("id")
	if presenceIdStr == "" {
		h.log.Error("[ArrivalPresence] presenceId not found in param")
		return errx.NotFound("presenceId not found in param")
	}

	presenceId, err := strconv.ParseUint(presenceIdStr, 10, 64)
	if err != nil {
		h.log.Error("[ArrivalPresence] failed to parse presenceId", zap.Error(err))
		return errx.BadRequest("invalid presence id param")
	}

	staffPresence, err := h.service.ArrivalPresence(presenceId, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[ArrivalPresence] failed to arrival presence", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, staffPresence, "success arrival presence")
}

func (h *PresenceHandler) DeparturePresence(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[DeparturePresence] userId not found in context")
		return errx.NotFound("userId not found in context")
	}

	presenceIdStr := c.Params("id")
	if presenceIdStr == "" {
		h.log.Error("[DeparturePresence] presenceId not found in param")
		return errx.NotFound("presenceId not found in param")
	}

	presenceId, err := strconv.ParseUint(presenceIdStr, 10, 64)
	if err != nil {
		h.log.Error("[DeparturePresence] failed to parse presenceId", zap.Error(err))
		return errx.BadRequest("invalid presence id param")
	}

	staffPresence, err := h.service.DeparturePresence(presenceId, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[DeparturePresence] failed to departure presence", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, staffPresence, "success departure presence")
}
