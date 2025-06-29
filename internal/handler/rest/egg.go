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

type EggHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IEggService
}

func (h *EggHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/eggs")
	v1.Post("/monitorings", middleware.Authentication(), h.CreateEggMonitoring)
	v1.Get("/monitorings", middleware.Authentication(), h.GetEggMonitorings)
	v1.Get("/monitorings/:id", middleware.Authentication(), h.GetEggMonitoringById)
	v1.Put("/monitorings/:id", middleware.Authentication(), h.UpdateEggMonitoring)
	v1.Delete("/monitorings/:id", middleware.Authentication(), h.DeleteEggMonitoring)
	v1.Patch("/monitorings/:id/takes", middleware.Authentication(), h.TakeEggMonitoring)

	v1.Get("/overview", middleware.Authentication(), h.GetEggOverview)
}

func NewEggHandler(log *zap.Logger, service service.IEggService, validator *validator.Validate) *EggHandler {
	return &EggHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *EggHandler) CreateEggMonitoring(c *fiber.Ctx) error {
	var request dto.CreateEggMonitoringRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	res, err := h.service.CreateEggMonitoring(request, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("failed to create egg monitoring", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create egg monitoring")
}

func (h *EggHandler) GetEggMonitorings(c *fiber.Ctx) error {
	var filter dto.GetEggMonitoringFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	res, err := h.service.GetEggMonitorings(filter)
	if err != nil {
		h.log.Error("failed to get egg monitorings", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get egg monitorings")
}

func (h *EggHandler) GetEggMonitoringById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	res, err := h.service.GetEggMonitoringById(id)
	if err != nil {
		h.log.Error("failed to get egg monitoring by id", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get egg monitoring by id")
}

func (h *EggHandler) UpdateEggMonitoring(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	var request dto.UpdateEggMonitoringRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	res, err := h.service.UpdateEggMonitoring(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update egg monitoring")
}

func (h *EggHandler) DeleteEggMonitoring(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	if err := h.service.DeleteEggMonitoring(id); err != nil {
		h.log.Error("failed to delete egg monitoring", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}

func (h *EggHandler) TakeEggMonitoring(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	res, err := h.service.TakeEggMonitoring(id, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("failed to take egg monitoring", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success take egg monitoring")
}

func (h *EggHandler) GetEggOverview(c *fiber.Ctx) error {
	var filter dto.GetEggOverviewFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	res, err := h.service.GetOverviewEggMonitoring(filter)
	if err != nil {
		h.log.Error("failed to get egg overview", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get egg overview")
}
