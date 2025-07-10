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

type ChickenHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IChickenService
}

func (h *ChickenHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/chickens")
	v1.Get("/overview", middleware.Authentication(), h.GetChickenOverview)

	v1.Post("/monitorings", middleware.Authentication(), h.CreateChickenMonitoring)
	v1.Get("/monitorings", middleware.Authentication(), h.GetChickenMonitorings)
	v1.Put("/monitorings/:id", middleware.Authentication(), h.UpdateChickenMonitoring)
	v1.Get("/monitorings/:id", middleware.Authentication(), h.GetChickenMonitoringById)
	v1.Delete("/monitorings/:id", middleware.Authentication(), h.DeleteChickenMonitoring)

	v1.Post("/healths/items", middleware.Authentication(), h.CreateChickenHealthItem)
	v1.Get("/healths/items", middleware.Authentication(), h.GetChickenHealthItems)
	v1.Get("/healths/items/:id", middleware.Authentication(), h.GetChickenHealthItemById)
	v1.Put("/healths/items/:id", middleware.Authentication(), h.UpdateChickenHealthItem)
	v1.Delete("/healths/items/:id", middleware.Authentication(), h.DeleteChickenHealthItem)

	v1.Post("/healths/monitorings", middleware.Authentication(), h.CreateChickenHealthMonitoring)
	v1.Get("/healths/monitorings/details/:chickenCageId", middleware.Authentication(), h.GetChickenHealthMonitoringDetails)
	v1.Get("/healths/monitorings/:id", middleware.Authentication(), h.GetChickenHealthMonitoringById)
	v1.Put("/healths/monitorings/:id", middleware.Authentication(), h.UpdateChickenHealthMonitoring)
	v1.Delete("/healths/monitorings/:id", middleware.Authentication(), h.DeleteChickenHealthMonitoring)
}

func NewChickenHandler(log *zap.Logger, service service.IChickenService, validator *validator.Validate) *ChickenHandler {
	return &ChickenHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *ChickenHandler) CreateChickenMonitoring(c *fiber.Ctx) error {
	var request dto.CreateChickenMonitoringRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("userId not found in context")
		return errx.Unauthorized("userId not found in context")
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("validation failed", zap.Error(err))
		return err
	}

	res, err := h.service.CreateChickenMonitoring(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create chicken monitoring")
}

func (h *ChickenHandler) GetChickenMonitorings(c *fiber.Ctx) error {
	var filter dto.GetChickenMonitoringFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	res, err := h.service.GetChickenMonitorings(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get chicken monitorings")
}

func (h *ChickenHandler) GetChickenMonitoringById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	res, err := h.service.GetChickenMonitoringById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		res,
		"success get chicken monitoring",
	)
}

func (h *ChickenHandler) UpdateChickenMonitoring(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	var request dto.UpdateChickenMonitoringRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("validation failed", zap.Error(err))
		return err
	}

	res, err := h.service.UpdateChickenMonitoring(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		res,
		"success update chicken monitoring",
	)
}

func (h *ChickenHandler) DeleteChickenMonitoring(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	err = h.service.DeleteChickenMonitoring(id)
	if err != nil {
		h.log.Error("failed to delete chicken monitoring", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}

func (h *ChickenHandler) CreateChickenHealthItem(c *fiber.Ctx) error {
	var request dto.CreateChickenHealthItemRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validaton error", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.CreateChickenHealthItem(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, data, "success create chicken health item")
}

func (h *ChickenHandler) GetChickenHealthItemById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	data, err := h.service.GetChickenHealthItemById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success to get chicken health item by id")
}

func (h *ChickenHandler) GetChickenHealthItems(c *fiber.Ctx) error {
	var filter dto.GetChickenHealthItemFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	data, err := h.service.GetChickenHealthItems(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success to get chicken health items")
}

func (h *ChickenHandler) UpdateChickenHealthItem(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	var request dto.UpdateChickenHealthItemRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.UpdateChickenHealthItem(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success update chicken health item")
}

func (h *ChickenHandler) DeleteChickenHealthItem(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id param", zap.Error(err))
		return err
	}

	err = h.service.DeleteChickenHealthItem(id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *ChickenHandler) CreateChickenHealthMonitoring(c *fiber.Ctx) error {
	var request dto.CreateChickenHealthMonitoringRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Warn("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.CreateChickenHealthMonitoring(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, data, "success create chicken health monitoring")
}

func (h *ChickenHandler) GetChickenHealthMonitoringDetails(c *fiber.Ctx) error {
	chickenCageId, err := strconv.ParseUint(c.Params("chickenCageId"), 10, 64)
	if err != nil {
		h.log.Error("invalid chicken cage id param", zap.Error(err))
		return errx.BadRequest("invalid chicken cage id param")
	}

	data, err := h.service.GetChickenHealthMonitoringDetails(chickenCageId)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get chicken health monitoring details")
}

func (h *ChickenHandler) GetChickenHealthMonitoringById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	data, err := h.service.GetChickenHealthMonitoringById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get chicken health monitoring by id")
}

func (h *ChickenHandler) UpdateChickenHealthMonitoring(c *fiber.Ctx) error {
	var request dto.UpdateChickenHealthMonitoringRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation failed", zap.Error(err))
		return err
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Warn("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.UpdateChickenHealthMonitoring(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success chicken health monitoring")
}

func (h *ChickenHandler) DeleteChickenHealthMonitoring(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	err = h.service.DeleteChickenHealthMonitoring(id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *ChickenHandler) GetChickenOverview(c *fiber.Ctx) error {
	var filter dto.GetChickenOverviewFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query param", zap.Error(err))
		return err
	}

	data, err := h.service.GetChickenOverview(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get chicken overview")
}
