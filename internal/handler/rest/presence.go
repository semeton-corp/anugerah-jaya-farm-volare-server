package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/service"
	"go.uber.org/zap"
)

type PresenceHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IPresenceService
}

func (a *PresenceHandler) SetEndpoint(router *fiber.App) {
	_ = router.Group("api/v1/presences")
}

func NewPresenceHandler(log *zap.Logger, service service.IPresenceService, validator *validator.Validate) *PresenceHandler {
	return &PresenceHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}
