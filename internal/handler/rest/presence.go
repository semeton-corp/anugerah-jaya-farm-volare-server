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

	v1.Get("/current/me", middleware.Authentication(), h.GetCurrentUserPresence)
	v1.Get("/me", middleware.Authentication(), h.GetCurrentUserPresences)
	v1.Patch("/:id", middleware.Authentication(), h.UpdateUserPresence)

	v1.Get("/locations/summaries", middleware.Authentication(), h.GetLocationPresenceSummaries)
	v1.Get("/users/summaries", middleware.Authentication(), h.GetLocationPresenceSummaries)
	v1.Get("/users/works/summaries", middleware.Authentication(), h.GetUserPresenceWorkDetailSummaries)
}

func NewPresenceHandler(log *zap.Logger, service service.IPresenceService, validator *validator.Validate) *PresenceHandler {
	return &PresenceHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *PresenceHandler) GetCurrentUserPresence(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found in context")
		return errx.NotFound("user id not found in context")
	}

	userPresence, err := h.service.GetCurrentUserPresence(uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, userPresence, "success get current user presence")
}

func (h *PresenceHandler) GetCurrentUserPresences(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found in context")
		return errx.NotFound("user id not found in context")
	}

	var filter dto.GetPresenceFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parsing query filter", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("failed to validate filter", zap.Error(err))
		return err
	}

	userPresences, err := h.service.GetUserPresencesByUserId(uuid.MustParse(userId), filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, userPresences, "success get all user presences")
}

func (h *PresenceHandler) UpdateUserPresence(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("userId not found in context")
		return errx.NotFound("userId not found in context")
	}

	presenceId, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse presenceId", zap.Error(err))
		return errx.BadRequest("invalid presence id param")
	}

	var request dto.UpdateUserPresenceRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	userPresence, err := h.service.UpdateUserPresence(presenceId, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, userPresence, "success update presence")
}

func (s *PresenceHandler) GetLocationPresenceSummaries(c *fiber.Ctx) error {
	data, err := s.service.GetLocationPresenceSummaries()
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get location presence summary")
}

func (s *PresenceHandler) GetUserPresenceSummaries(c *fiber.Ctx) error {
	var filter dto.GetUserPresenceSummaryFilter
	if err := c.QueryParser(&filter); err != nil {
		s.log.Error("failed to parse query param", zap.Error(err))
		return err
	}

	if err := s.validator.Struct(&filter); err != nil {
		s.log.Error("validation error", zap.Error(err))
		return err
	}

	data, err := s.service.GetUserPresenceSummaries(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get user presence summaries")
}

func (s *PresenceHandler) GetUserPresenceWorkDetailSummaries(c *fiber.Ctx) error {
	var filter dto.GetUserPresenceWorkDetailSummaryFilter
	if err := c.QueryParser(&filter); err != nil {
		s.log.Error("failed to parse query param", zap.Error(err))
		return err
	}

	if err := s.validator.Struct(&filter); err != nil {
		s.log.Error("validation error", zap.Error(err))
		return err
	}

	data, err := s.service.GetUserPresenceWorkDetailSummaries(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get user presense work detail summaries")
}
