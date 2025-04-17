package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/middleware"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/service"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/response"
	"go.uber.org/zap"
)

type RoleHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IRoleService
}

func (h *RoleHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/roles")
	v1.Get("/", middleware.Authentication(), h.GetRoles)
}

func NewRoleHandler(log *zap.Logger, service service.IRoleService, validator *validator.Validate) *RoleHandler {
	return &RoleHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *RoleHandler) GetRoles(c *fiber.Ctx) error {
	roles, err := h.service.GetRoles()
	if err != nil {
		h.log.Error("[GetRoles] failed to get roles", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		roles,
		"success get roles",
	)
}
