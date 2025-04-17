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

func (a *AuthenticationHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/authentication")
	v1.Post("/signin", a.SignIn)
	v1.Post("/forgot-password", a.ForgotPassword)

	v1.Post("/signup", middleware.Authentication(), a.SignUp)
	v1.Post("/change-password", middleware.Authentication(), a.ChangePassword)
	v1.Put("/:id", middleware.Authentication(), a.UpdateAccount)
	v1.Delete("/:id", middleware.Authentication(), a.DeleteAccount)
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

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[SignUp] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	res, err := a.service.SignUp(request, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[SignUp] failed to sign up", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusCreated,
		res,
		"success sign up",
	)
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

	res, err := a.service.SignIn(request)
	if err != nil {
		a.log.Error("[SignIn] failed to sign in", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		res,
		"success sign in",
	)
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

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[ForgotPassword] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	res, err := a.service.ForgotPassword(request, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[ForgotPassword] failed to forgot password", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		res,
		"success reset password",
	)
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

	res, err := a.service.ChangePassword(request, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[ChangePassword] failed to reset password", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		res,
		"success change password",
	)
}

func (a *AuthenticationHandler) UpdateAccount(c *fiber.Ctx) error {
	var request dto.UpdateAccountRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[UpdateAccount] failed to parse request", zap.Error(err))
		return err
	}

	idParam := c.Params("id")
	if idParam == "" {
		a.log.Error("[UpdateAccount] id is required")
		return errx.BadRequest("id is required")
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[UpdateAccount] failed to validate request", zap.Error(err))
		return err
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[UpdateAccount] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	res, err := a.service.UpdateAccount(uuid.MustParse(idParam), request, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[UpdateAccount] failed to update account", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		res,
		"success update account",
	)
}

func (a *AuthenticationHandler) DeleteAccount(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		a.log.Error("[DeleteAccount] id is required")
		return errx.BadRequest("id is required")
	}

	if err := a.service.DeleteAccount(uuid.MustParse(idParam)); err != nil {
		a.log.Error("[DeleteAccount] failed to delete account", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}
