package rest

import (
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

type NotificationHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.INotificationService
}

func (h *NotificationHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/notifications")
	v1.Get("", middleware.Authentication(), h.GetNotifications)
	v1.Patch("", middleware.Authentication(), h.MarkNotifications)
}

func NewNotificationHandler(log *zap.Logger, service service.INotificationService, validator *validator.Validate) *NotificationHandler {
	return &NotificationHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	var filter dto.GetNotificationFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed parse query param", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&filter); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return err
	}

	data, err := h.service.GetNotifications(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get notifications")
}

func (h *NotificationHandler) MarkNotifications(c *fiber.Ctx) error {
	var request dto.MarkNotificationRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed parse body request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.MarkNotifications(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success marked notifications as read")
}
