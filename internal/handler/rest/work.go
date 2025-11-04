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

type WorkHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IWorkService
}

func (h *WorkHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/works")
	v1.Get("/overview", middleware.Authentication(), h.GetWorkOverview)

	v1.Get("/me", middleware.Authentication(), h.GetSelfWorkUser)

	v1.Post("/dailies", middleware.Authentication(), h.SaveDailyWorks)
	v1.Get("/dailies/summaries", middleware.Authentication(), h.GetDailyWorksSummariesBasedOnRole)
	v1.Get("/dailies/:roleId", middleware.Authentication(), h.GetDailyWorksByRoleId)
	v1.Put("/dailies/users/:id", middleware.Authentication(), h.UpdateDailyWorkUser)
	v1.Delete("/dailies/:id", middleware.Authentication(), h.DeleteDailyWork)

	v1.Put("/additionals/users/:id", middleware.Authentication(), h.UpdateAdditionalWorkUser)
	v1.Delete("/additionals/users/:id", middleware.Authentication(), h.DeleteAdditionalWorkUser)
	v1.Post("/additionals", middleware.Authentication(), h.CreateAdditionalWork)
	v1.Get("/additionals", middleware.Authentication(), h.GetAdditionalWorks)
	v1.Get("/additionals/:id", middleware.Authentication(), h.GetAdditionalWorkById)
	v1.Put("/additionals/:id", middleware.Authentication(), h.UpdateAdditionalWork)
	v1.Delete("/additionals/:id", middleware.Authentication(), h.DeleteAdditionalWork)
	v1.Post("/additionals/takes/:id", middleware.Authentication(), h.TakeAdditionalWork)

	v1.Get("/dailies/users/:userId", middleware.Authentication(), h.GetDailyWorkUsersByUserId)
	v1.Get("/additionals/users/:userId", middleware.Authentication(), h.GetAdditionalWorkUsersByUserId)
}

func NewWorkHandler(log *zap.Logger, service service.IWorkService, validator *validator.Validate) *WorkHandler {
	return &WorkHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *WorkHandler) GetWorkOverview(c *fiber.Ctx) error {
	res, err := h.service.GetWorkOverview()
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get work overview")
}

func (h *WorkHandler) SaveDailyWorks(c *fiber.Ctx) error {
	var request dto.CreateDailyWorkRequest
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
		h.log.Warn("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	res, err := h.service.SaveDailyWorks(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success save daily works")
}

func (h *WorkHandler) GetDailyWorksSummariesBasedOnRole(c *fiber.Ctx) error {
	res, err := h.service.GetDailyWorkSummariesBasedOnRole()
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get daily works based on role")
}

func (h *WorkHandler) GetDailyWorksByRoleId(c *fiber.Ctx) error {
	roleId, err := strconv.ParseUint(c.Params("roleId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse role id", zap.Error(err))
		return errx.BadRequest("role id must be a number")
	}

	res, err := h.service.GetDailyWorksByRoleId(roleId)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get daily works by role id")
}

func (h *WorkHandler) CreateAdditionalWork(c *fiber.Ctx) error {
	var request dto.CreateAdditionalWorkRequest
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
		h.log.Warn("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	res, err := h.service.CreateAdditionalWork(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create additional work")
}

func (h *WorkHandler) GetAdditionalWorks(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Warn("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	var filter dto.GetAdditonalWorkFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&filter); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return errx.BadRequest("error validation")
	}

	res, err := h.service.GetAdditionalWorks(filter, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get additional works")
}

func (h *WorkHandler) GetAdditionalWorkById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse work id", zap.Error(err))
		return errx.BadRequest("work id must be a number")
	}

	res, err := h.service.GetAdditionalWorkById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get additional work by id")
}

func (h *WorkHandler) UpdateAdditionalWork(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse work id", zap.Error(err))
		return errx.BadRequest("work id must be a number")
	}

	var request dto.UpdateAdditionalWorkRequest
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
		h.log.Warn("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	res, err := h.service.UpdateAdditionalWork(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update additional work")
}

func (h *WorkHandler) DeleteAdditionalWork(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return errx.BadRequest(" id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse work id", zap.Error(err))
		return errx.BadRequest("work id must be a number")
	}

	err = h.service.DeleteAdditionalWork(id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *WorkHandler) GetSelfWorkUser(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Warn("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	res, err := h.service.GetUserWorksByUserId(uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get own user work")
}

func (h *WorkHandler) TakeAdditionalWork(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Warn("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse work id", zap.Error(err))
		return errx.BadRequest("work id must be a number")
	}

	resp, err := h.service.TakeAdditionalWork(id, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, resp, "success take additional work")
}

func (h *WorkHandler) UpdateAdditionalWorkUser(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Warn("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}
	idParam := c.Params("id")
	if idParam == "" {
		return errx.BadRequest(" id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse work id", zap.Error(err))
		return errx.BadRequest("work id must be a number")
	}

	var request dto.UpdateAdditionalWorkUserRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	res, err := h.service.UpdateAdditionalWorkUser(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update additional work user")
}

func (h *WorkHandler) UpdateDailyWorkUser(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Warn("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	IdParam := c.Params("id")
	if IdParam == "" {
		return errx.BadRequest(" id is required")
	}

	Id, err := strconv.ParseUint(IdParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse work id", zap.Error(err))
		return errx.BadRequest("work id must be a number")
	}

	var request dto.UpdateDailyWorkUserRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	res, err := h.service.UpdateDailyWorkUser(Id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update daily work user")
}

func (h *WorkHandler) DeleteDailyWork(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return errx.BadRequest(" id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse work id", zap.Error(err))
		return errx.BadRequest("work id must be a number")
	}

	err = h.service.DeleteDailyWork(id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *WorkHandler) DeleteAdditionalWorkUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return err
	}

	err = h.service.DeleteAdditionalWorkUser(id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *WorkHandler) GetDailyWorkUsersByUserId(c *fiber.Ctx) error {
	var filter dto.GetDailyWorkUserFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&filter); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	data, err := h.service.GetDailyWorkUserByUserId(uuid.MustParse(c.Params("userId")), filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get daily work user by user id")
}

func (h *WorkHandler) GetAdditionalWorkUsersByUserId(c *fiber.Ctx) error {
	var filter dto.GetAdditionalWorkUserFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&filter); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	data, err := h.service.GetAdditionalWorkUserByUserId(uuid.MustParse(c.Params("userId")), filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get additional work user by user id")
}
