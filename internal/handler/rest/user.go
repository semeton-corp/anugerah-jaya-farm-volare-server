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

	v1.Put("/me", middleware.Authentication(), h.UpdateSelfUser)
	v1.Get("/me", middleware.Authentication(), h.GetSelfUser)

	v1.Get("/overview", middleware.Authentication(), h.GetUserOverviews)
	v1.Get("/overview/:id", middleware.Authentication(), h.GetOverviewUser)

	v1.Get("", middleware.Authentication(), h.GetUsers)
	v1.Get("/:id", middleware.Authentication(), h.GetUserById)
	v1.Put("/:id", middleware.Authentication(), h.UpdateUser)

	v1.Get("/performances/overview", middleware.Authentication(), h.GetUserPerformanceOverview)
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
		h.log.Error("id param not found")
		return errx.BadRequest("id param not found")
	}

	resp, err := h.service.GetUserById(uuid.MustParse(idParam))
	if err != nil {
		h.log.Error("failed to get user by id")
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get user by id")
}

func (h *UserHandler) GetSelfUser(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found in context")
		return errx.BadRequest("user id not found in context")
	}

	resp, err := h.service.GetUserById(uuid.MustParse(userId))
	if err != nil {
		h.log.Error("failed to get user by id")
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get user by id")
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("id param not found")
		return errx.BadRequest("id param not found")
	}

	var request dto.UpdateUserRequest
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

	res, err := h.service.UpdateUser(uuid.MustParse(idParam), request, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("failed to update account", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update account")
}

func (h *UserHandler) UpdateSelfUser(c *fiber.Ctx) error {
	var request dto.UpdateUserRequest
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

	res, err := h.service.UpdateUser(uuid.MustParse(userId), request, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("failed to update account", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update account")
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	var filter dto.GetUserListFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parsing query filter", zap.Error(err))
		return err
	}

	resp, err := h.service.GetUsers(filter)
	if err != nil {
		h.log.Error("failed to get users")
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success to get users")
}

func (h *UserHandler) GetOverviewUser(c *fiber.Ctx) error {
	var filter dto.GetUserOverviewFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("failed to validate query", zap.Error(err))
		return err
	}

	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("id param is empty")
		return errx.BadRequest("id param is empty")
	}

	resp, err := h.service.GetUserOverview(uuid.MustParse(idParam), filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success to get overview user")
}

func (h *UserHandler) GetUserOverviews(c *fiber.Ctx) error {
	var filter dto.GetUserOverviewListFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query param", zap.Error(err))
		return err
	}

	data, err := h.service.GetUserOverviewList(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get user overview list")
}

func (h *UserHandler) GetUserPerformanceOverview(c *fiber.Ctx) error {
	var filter dto.GetUserPerformanceOverviewFilter

	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed parse query", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&filter); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return err
	}

	data, err := h.service.GetUserPerformanceOverview(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get user performance overview")
}
