package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/service"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"go.uber.org/zap"
)

type AuthenticationHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IAuthenticationService
}

func (a *AuthenticationHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/authentication")
	v1.Post("/signup", a.SignUp)
	v1.Post("/signin", a.SignIn)
	v1.Post("/forgot-password", a.ForgotPassword)

	// this one need middleware
	v1.Post("/change-password", a.ChangePassword)
}

func NewAuthenticationHandler(log *zap.Logger, service service.IAuthenticationService, validator *validator.Validate) *AuthenticationHandler {
	return &AuthenticationHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (a *AuthenticationHandler) SignUp(c *fiber.Ctx) error {
	var request dto.SignUpRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[SignUp] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[SignUp] failed to validate request", zap.Error(err))
		return err
	}

	response, err := a.service.SignUp(request)
	if err != nil {
		a.log.Error("[SignUp] failed to sign up", zap.Error(err))
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"data":    response,
		"message": "success create account",
	})
}

func (a *AuthenticationHandler) SignIn(c *fiber.Ctx) error {
	var request dto.SignInRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[SignIn] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[SignIn] failed to validate request", zap.Error(err))
		return err
	}

	response, err := a.service.SignIn(request)
	if err != nil {
		a.log.Error("[SignIn] failed to sign in", zap.Error(err))
		return err
	}

	return c.Status(fiber.StatusOK).Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": fiber.StatusOK,
		"data":   response,
	})
}

func (a *AuthenticationHandler) ForgotPassword(c *fiber.Ctx) error {
	var request dto.ForgotPasswordRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[ForgotPassword] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[ForgotPassword] failed to validate request", zap.Error(err))
		return err
	}

	response, err := a.service.ForgotPassword(request)
	if err != nil {
		a.log.Error("[ForgotPassword] failed to forgot password", zap.Error(err))
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (a *AuthenticationHandler) ChangePassword(c *fiber.Ctx) error {
	var request dto.ChangePasswordRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[ChangePassword] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[ChangePassword] failed to validate request", zap.Error(err))
		return err
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[ChangePassword] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	response, err := a.service.ChangePassword(request, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[ChangePassword] failed to reset password", zap.Error(err))
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
