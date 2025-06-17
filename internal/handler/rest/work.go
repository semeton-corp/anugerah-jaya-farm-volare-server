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
	v1.Get("/me", middleware.Authentication(), h.GetOwnStaffWork)
	v1.Post("/dailies", middleware.Authentication(), h.CreateAndUpdateDailyWorks)
	v1.Get("/dailies", middleware.Authentication(), h.GetDailyWorksBasedOnRole)
	v1.Get("/dailies/:roleId", middleware.Authentication(), h.GetDailyWorksByRoleId)
	v1.Put("/dailies/staffs/:id", middleware.Authentication(), h.UpdateDailyWorkStaff)
	v1.Delete("/dailies/:id", middleware.Authentication(), h.DeleteDailyWork)

	v1.Put("/additionals/staffs/:id", middleware.Authentication(), h.UpdateAdditionalWorkStaff)
	v1.Post("/additionals", middleware.Authentication(), h.CreateAdditionalWork)
	v1.Get("/additionals", middleware.Authentication(), h.GetAdditionalWorks)
	v1.Get("/additionals/:id", middleware.Authentication(), h.GetAdditionalWorkById)
	v1.Put("/additionals/:id", middleware.Authentication(), h.UpdateAdditionalWork)
	v1.Delete("/additionals/:id", middleware.Authentication(), h.DeleteAdditionalWork)
	v1.Post("/additionals/takes/:id", middleware.Authentication(), h.TakeAdditionalWork)
}

func NewWorkHandler(log *zap.Logger, service service.IWorkService, validator *validator.Validate) *WorkHandler {
	return &WorkHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *WorkHandler) CreateAndUpdateDailyWorks(c *fiber.Ctx) error {
	var request dto.CreateDailyWorkRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	idCtx := c.Locals("accountId").(string)

	res, err := h.service.CreateAndUpdateDailyWorks(request, uuid.MustParse(idCtx))
	if err != nil {
		h.log.Error("failed to create and update daily works", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create and update daily works")
}

func (h *WorkHandler) GetDailyWorksBasedOnRole(c *fiber.Ctx) error {
	var filter dto.GetDailyWorkBasedOnRoleFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	res, err := h.service.GetDailyWorksBasedOnRole(filter)
	if err != nil {
		h.log.Error("failed to get daily works based on role", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get daily works based on role")
}

func (h *WorkHandler) GetDailyWorksByRoleId(c *fiber.Ctx) error {
	roleIdParam := c.Params("roleId")
	if roleIdParam == "" {
		return errx.BadRequest("role id is required")
	}

	roleId, err := strconv.ParseUint(roleIdParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse role id", zap.Error(err))
		return errx.BadRequest("role id must be a number")
	}

	res, err := h.service.GetDailyWorksByRoleId(roleId)
	if err != nil {
		h.log.Error("failed to get daily works by role id", zap.Error(err))
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

	idCtx := c.Locals("accountId").(string)

	res, err := h.service.CreateAdditionalWork(request, uuid.MustParse(idCtx))
	if err != nil {
		h.log.Error("failed to create additional work", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create additional work")
}

func (h *WorkHandler) GetAdditionalWorks(c *fiber.Ctx) error {
	var filter dto.GetAdditonalWorkFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	res, err := h.service.GetAdditionalWorks(filter)
	if err != nil {
		h.log.Error("failed to get additional works", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get additional works")
}

func (h *WorkHandler) GetAdditionalWorkById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return errx.BadRequest(" id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse work id", zap.Error(err))
		return errx.BadRequest("work id must be a number")
	}

	res, err := h.service.GetAdditionalWorkById(id)
	if err != nil {
		h.log.Error("failed to get additional work by id", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get additional work by id")
}

func (h *WorkHandler) UpdateAdditionalWork(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return errx.BadRequest(" id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
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

	idCtx := c.Locals("accountId").(string)

	res, err := h.service.UpdateAdditionalWork(id, request, uuid.MustParse(idCtx))
	if err != nil {
		h.log.Error("failed to update additional work", zap.Error(err))
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
		h.log.Error("failed to delete additional work", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}

func (h *WorkHandler) GetOwnStaffWork(c *fiber.Ctx) error {
	idCtx := c.Locals("accountId").(string)

	res, err := h.service.GetStaffWorksByStaffId(uuid.MustParse(idCtx))
	if err != nil {
		h.log.Error("failed to get own staff work", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get own staff work")
}

func (h *WorkHandler) TakeAdditionalWork(c *fiber.Ctx) error {
	accountIdCtx := c.Locals("accountId").(string)

	idParam := c.Params("id")
	if idParam == "" {
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse work id", zap.Error(err))
		return errx.BadRequest("work id must be a number")
	}

	resp, err := h.service.TakeAdditionalWork(id, uuid.MustParse(accountIdCtx))
	if err != nil {
		h.log.Error("failed to take additional work", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, resp, "success take additional work")
}

func (h *WorkHandler) UpdateAdditionalWorkStaff(c *fiber.Ctx) error {
	accountIdCtx := c.Locals("accountId").(string)

	idParam := c.Params("id")
	if idParam == "" {
		return errx.BadRequest(" id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse work id", zap.Error(err))
		return errx.BadRequest("work id must be a number")
	}

	var request dto.UpdateAdditionalWorkStaffRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	res, err := h.service.UpdateAdditionalWorkStaff(id, request, uuid.MustParse(accountIdCtx))
	if err != nil {
		h.log.Error("failed to update additional work staff", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update additional work staff")
}

func (h *WorkHandler) UpdateDailyWorkStaff(c *fiber.Ctx) error {
	accountIdCtx := c.Locals("accountId").(string)

	IdParam := c.Params("id")
	if IdParam == "" {
		return errx.BadRequest(" id is required")
	}

	Id, err := strconv.ParseUint(IdParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse work id", zap.Error(err))
		return errx.BadRequest("work id must be a number")
	}

	var request dto.UpdateDailyWorkStaffRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	res, err := h.service.UpdateDailyWorkStaff(Id, request, uuid.MustParse(accountIdCtx))
	if err != nil {
		h.log.Error("failed to update daily work staff", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update daily work staff")
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
		h.log.Error("failed to delete daily work", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}
