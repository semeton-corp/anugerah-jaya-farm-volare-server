package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/service"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/response"
	"go.uber.org/zap"
)

type LocationHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.ILocationService
}

func (h *LocationHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/locations")
	v1.Get("/", h.GetLocations)
}

func NewLocationHandler(log *zap.Logger, service service.ILocationService, validator *validator.Validate) *LocationHandler {
	return &LocationHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *LocationHandler) GetLocations(c *fiber.Ctx) error {
	data, err := h.service.GetLocations()
	if err != nil {
		h.log.Error("[GetLocations] failed to get locations", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get locations")
}
