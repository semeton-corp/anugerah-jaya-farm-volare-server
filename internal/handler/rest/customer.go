package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/service"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/response"
	"go.uber.org/zap"
)

type CustomerHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.ICustomerService
}

func (h *CustomerHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/customers")
	v1.Get("/", h.GetCustomers)
}

func NewCustomerHandler(log *zap.Logger, service service.ICustomerService, validator *validator.Validate) *CustomerHandler {
	return &CustomerHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *CustomerHandler) GetCustomers(c *fiber.Ctx) error {
	data, err := h.service.GetCustomers()
	if err != nil {
		h.log.Error("failed to get customers", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get customers")
}
