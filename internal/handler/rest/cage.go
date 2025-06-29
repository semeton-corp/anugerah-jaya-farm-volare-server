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

type CageHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.ICageService
}

func (h *CageHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/cages")
	v1.Get("/", middleware.Authentication(), h.GetCages)
	v1.Post("/", middleware.Authentication(), h.CreateCage)
	v1.Put("/:id", middleware.Authentication(), h.UpdateCage)
	v1.Delete("/:id", middleware.Authentication(), h.DeleteCage)

	v1.Get("/chickens", middleware.Authentication(), h.GetChickenCages)
	v1.Get("/chickens/:id", middleware.Authentication(), h.GetChickenCageById)
}

func NewCageHandler(log *zap.Logger, service service.ICageService, validator *validator.Validate) *CageHandler {
	return &CageHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *CageHandler) GetCages(c *fiber.Ctx) error {
	var filter dto.GetCageFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query filter", zap.Error(err))
		return err
	}

	res, err := h.service.GetCages(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get cages")
}

func (h *CageHandler) CreateCage(c *fiber.Ctx) error {
	var request dto.CreateCageRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation failed", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("userId not found in context")
		return errx.Unauthorized("userId not found in context")
	}

	data, err := h.service.CreateCage(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, data, "success create cage")
}

func (h *CageHandler) UpdateCage(c *fiber.Ctx) error {
	var request dto.UpdateCageRequest
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
		h.log.Error("failed parse id param")
		return errx.BadRequest("invalid id param")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("userId not found in context")
		return errx.Unauthorized("userId not found in context")
	}

	data, err := h.service.UpdateCage(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success update cage")
}

func (h *CageHandler) DeleteCage(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed parse id param")
		return errx.BadRequest("invalid id param")
	}

	err = h.service.DeleteCage(id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *CageHandler) GetChickenCages(c *fiber.Ctx) error {
	var filter dto.GetChickenCageFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query filter", zap.Error(err))
		return err
	}

	res, err := h.service.GetChickenCages(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get chicken cages")
}

func (h *CageHandler) GetChickenCageById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Warn("failed to parse id param")
		return errx.BadRequest("invalid id param")
	}

	res, err := h.service.GetChickenCageById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get chicken cage by id")
}
