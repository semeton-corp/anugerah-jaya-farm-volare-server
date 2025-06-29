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

type AuthenticationHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IAuthenticationService
}

func (h *AuthenticationHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/authentication")
	v1.Post("/signin", h.SignIn)
	v1.Post("/forgot-password", h.ForgotPassword)

	v1.Post("/signup", middleware.Authentication(), h.SignUp)
	v1.Post("/change-password", middleware.Authentication(), h.ChangePassword)
	v1.Delete("/:id", middleware.Authentication(), h.DeleteUser)
}

func NewAuthenticationHandler(log *zap.Logger, service service.IAuthenticationService, validator *validator.Validate) *AuthenticationHandler {
	return &AuthenticationHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *AuthenticationHandler) SignUp(c *fiber.Ctx) error {
	var request dto.SignUpRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("validation failed", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("userId not found in context")
		return errx.Unauthorized("no userId in context")
	}

	res, err := h.service.SignUp(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success sign up")
}

func (h *AuthenticationHandler) SignIn(c *fiber.Ctx) error {
	var request dto.SignInRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[SignIn] failed to validate request", zap.Error(err))
		return err
	}

	res, err := h.service.SignIn(request)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success sign in")
}

func (h *AuthenticationHandler) ForgotPassword(c *fiber.Ctx) error {
	var request dto.ForgotPasswordRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	res, err := h.service.ForgotPassword(request)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success reset password")
}

func (h *AuthenticationHandler) ChangePassword(c *fiber.Ctx) error {
	var request dto.ChangePasswordRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	res, err := h.service.ChangePassword(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success change password")
}

func (h *AuthenticationHandler) DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Warn("id in param is required")
		return errx.BadRequest("id is required")
	}

	if err := h.service.DeleteUser(uuid.MustParse(idParam)); err != nil {
		return err
	}

	return response.NoContentResponse(c)
}
