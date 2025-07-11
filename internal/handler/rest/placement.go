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

type PlacementHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IPlacementService
}

func (h *PlacementHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/placements")
	v1.Get("/stores/me", middleware.Authentication(), h.GetCurrentUserStorePlacement)
	v1.Post("/stores", middleware.Authentication(), h.CreateStorePlacement)
	v1.Delete("/stores/:storeId/", middleware.Authentication(), h.DeleteStorePlacementByUserId)

	v1.Get("/warehouses/me", middleware.Authentication(), h.GetCurrentUserWarehousePlacement)
	v1.Post("/warehouses", middleware.Authentication(), h.CreateWarehousePlacement)
	v1.Delete("/warehouses/:warehouseId/users/:userId", middleware.Authentication(), h.DeleteWarehousePlacementByUserId)

	v1.Get("/cages/me", middleware.Authentication(), h.GetCurrentUserCagePlacement)
	v1.Post("/cages", middleware.Authentication(), h.UpdateCagePlacement)
	v1.Delete("/cages/:cageId/users/:userId", middleware.Authentication(), h.DeleteCagePlacementByUserIdAndCageId)
}

func NewPlacementHandler(log *zap.Logger, service service.IPlacementService, validator *validator.Validate) *PlacementHandler {
	return &PlacementHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *PlacementHandler) GetCurrentUserStorePlacement(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Warn("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.GetStorePlacementByUserId(uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get store placement")
}

func (h *PlacementHandler) GetCurrentUserWarehousePlacement(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Warn("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.GetWarehousePlacementByUserId(uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get warehouse placement")
}

func (h *PlacementHandler) GetCurrentUserCagePlacement(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Warn("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.GetCagePlacementByUserId(uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get cage placement")
}

func (h *PlacementHandler) CreateStorePlacement(c *fiber.Ctx) error {
	var request dto.CreateStorePlacementRequest
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

	data, err := h.service.CreateStorePlacement(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, data, "success create store placement")
}

func (h *PlacementHandler) CreateWarehousePlacement(c *fiber.Ctx) error {
	var request dto.CreateWarehousePlacementRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Warn("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Warn("user id not found context")
		return errx.Unauthorized("user id not founc in context")
	}

	data, err := h.service.CreateWarehousePlacement(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, data, "success create warehouse placement")
}

func (h *PlacementHandler) UpdateCagePlacement(c *fiber.Ctx) error {
	var requests []dto.UpdateCagePlacementRequest
	if err := c.BodyParser(&requests); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	for i := range requests {
		if err := h.validator.Struct(&requests[i]); err != nil {
			h.log.Error("validation error", zap.Error(err))
			return err
		}
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Warn("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.UpdateCagePlacement(requests, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success update cage placement")
}

func (h *PlacementHandler) DeleteStorePlacementByUserId(c *fiber.Ctx) error {
	err := h.service.DeleteStorePlacementByUserId(uuid.MustParse(c.Params("userId")))
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *PlacementHandler) DeleteWarehousePlacementByUserId(c *fiber.Ctx) error {
	err := h.service.DeleteWarehousePlacementByUserId(uuid.MustParse(c.Params("userId")))
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *PlacementHandler) DeleteCagePlacementByUserIdAndCageId(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("cageId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id param", zap.Error(err))
		return err
	}

	err = h.service.DeleteCagePlacementByUserIdAndCageId(uuid.MustParse(c.Params("userId")), id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}
