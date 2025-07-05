package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	v1.Get("/store/me", middleware.Authentication(), h.GetCurrentUserStorePlacement)
	v1.Get("/warehouse/me", middleware.Authentication(), h.GetCurrentUserWarehousePlacement)
	v1.Get("/cage/me", middleware.Authentication(), h.GetCurrentUserCagePlacement)
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
