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

type UserHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IUserService
}

func (h *UserHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/users")

	v1.Put("/me", middleware.Authentication(), h.UpdateOwnProfile)
	v1.Get("/me", middleware.Authentication(), h.GetOwnProfile)
	v1.Get("", middleware.Authentication(), h.GetUsers)
	v1.Get("/:id", middleware.Authentication(), h.GetUserById)
	v1.Put("/:id", middleware.Authentication(), h.UpdateUser)
	v1.Get("/overview/:id", middleware.Authentication(), h.GetOverviewUser)
}

func NewUserHandler(log *zap.Logger, service service.IUserService, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *UserHandler) GetUserById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[GetUserById] id param not found")
		return errx.BadRequest("id param not found")
	}

	resp, err := h.service.GetUserById(uuid.MustParse(idParam))
	if err != nil {
		h.log.Error("[GetUserById] failed to get staff by id")
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get staff by id")
}

func (h *UserHandler) GetOwnProfile(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("userId not found in context")
		return errx.BadRequest("userId not found in context")
	}

	resp, err := h.service.GetUserById(uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[GetStaffById] failed to get staff by id")
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get staff by id")
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[UpdateUser] id param not found")
		return errx.BadRequest("id param not found")
	}

	var request dto.UpdateUserRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[UpdateUser] failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[UpdateUser] failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[UpdateUser] failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	res, err := h.service.UpdateUser(uuid.MustParse(idParam), request, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[UpdateAccount] failed to update account", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update account")
}

func (h *UserHandler) UpdateOwnProfile(c *fiber.Ctx) error {
	var request dto.UpdateUserRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[UpdateAccount] failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[UpdateAccount] failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[UpdateAccount] failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	res, err := h.service.UpdateUser(uuid.MustParse(userId), request, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[UpdateAccount] failed to update account", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update account")
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	var filter dto.GetUserFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("[GetStaffs] failed to parsing query filter", zap.Error(err))
		return err
	}

	resp, err := h.service.GetUsers(filter)
	if err != nil {
		h.log.Error("[GetStaffs] failed to get staffs")
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success to get staffs")
}

func (h *UserHandler) GetOverviewUser(c *fiber.Ctx) error {
	var filter dto.GetUserOverviewFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("[GetOverviewStaff] failed to parse query", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("[GetOverviewStaff] failed to validate query", zap.Error(err))
		return err
	}

	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[GetOverviewStaff] id param is empty")
		return errx.BadRequest("id param is empty")
	}

	resp, err := h.service.GetOverviewUser(uuid.MustParse(idParam), filter)
	if err != nil {
		h.log.Error("[GetOverviewStaff] failed to get overview staff", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success to get overview staff")
}
