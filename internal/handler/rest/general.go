package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/middleware"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/service"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/response"
	"go.uber.org/zap"
)

type GeneralHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IGeneralService
}

func (h *GeneralHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/generals")
	v1.Get("/overview", middleware.Authentication(), h.GetGeneralOverview)
}

func NewGeneralHandler(log *zap.Logger, service service.IGeneralService, validator *validator.Validate) *GeneralHandler {
	return &GeneralHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *GeneralHandler) GetGeneralOverview(c *fiber.Ctx) error {
	data, err := h.service.GetGeneralOverview()
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get general overview")
}
