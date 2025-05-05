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

func (a *WorkHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/works")
	v1.Get("/me", middleware.Authentication(), a.GetOwnStaffWork)
	v1.Post("/dailies", middleware.Authentication(), a.CreateAndUpdateDailyWorks)
	v1.Get("/dailies", middleware.Authentication(), a.GetDailyWorksBasedOnRole)
	v1.Get("/dailies/:roleId", middleware.Authentication(), a.GetDailyWorksByRoleId)
	v1.Put("/dailies/staffs/:id", middleware.Authentication(), a.UpdateDailyWorkStaff)
	v1.Delete("/dailies/:id", middleware.Authentication(), a.DeleteDailyWork)

	v1.Put("/additionals/staffs/:id", middleware.Authentication(), a.UpdateAdditionalWorkStaff)
	v1.Post("/additionals", middleware.Authentication(), a.CreateAdditionalWork)
	v1.Get("/additionals", middleware.Authentication(), a.GetAdditionalWorks)
	v1.Get("/additionals/:id", middleware.Authentication(), a.GetAdditionalWorkById)
	v1.Put("/additionals/:id", middleware.Authentication(), a.UpdateAdditionalWork)
	v1.Delete("/additionals/:id", middleware.Authentication(), a.DeleteAdditionalWork)
	v1.Post("/additionals/takes/:id", middleware.Authentication(), a.TakeAdditionalWork)
}

func NewWorkHandler(log *zap.Logger, service service.IWorkService, validator *validator.Validate) *WorkHandler {
	return &WorkHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (a *WorkHandler) CreateAndUpdateDailyWorks(c *fiber.Ctx) error {
	var request dto.CreateDailyWorkRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	idCtx := c.Locals("accountId").(string)

	res, err := a.service.CreateAndUpdateDailyWorks(request, uuid.MustParse(idCtx))
	if err != nil {
		a.log.Error("failed to create and update daily works", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create and update daily works")
}

func (a *WorkHandler) GetDailyWorksBasedOnRole(c *fiber.Ctx) error {
	res, err := a.service.GetDailyWorksBasedOnRole()
	if err != nil {
		a.log.Error("failed to get daily works based on role", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get daily works based on role")
}

func (a *WorkHandler) GetDailyWorksByRoleId(c *fiber.Ctx) error {
	roleIdParam := c.Params("roleId")
	if roleIdParam == "" {
		return errx.BadRequest("role id is required")
	}

	roleId, err := strconv.ParseUint(roleIdParam, 10, 64)
	if err != nil {
		a.log.Error("failed to parse role id", zap.Error(err))
		return errx.BadRequest("role id must be a number")
	}

	res, err := a.service.GetDailyWorksByRoleId(roleId)
	if err != nil {
		a.log.Error("failed to get daily works by role id", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get daily works by role id")
}

func (a *WorkHandler) CreateAdditionalWork(c *fiber.Ctx) error {
	var request dto.CreateAdditionalWorkRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	idCtx := c.Locals("accountId").(string)

	res, err := a.service.CreateAdditionalWork(request, uuid.MustParse(idCtx))
	if err != nil {
		a.log.Error("failed to create additional work", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create additional work")
}

func (a *WorkHandler) GetAdditionalWorks(c *fiber.Ctx) error {
	var filter dto.GetAdditonalWorkFilter
	if err := c.QueryParser(&filter); err != nil {
		a.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	res, err := a.service.GetAdditionalWorks(filter)
	if err != nil {
		a.log.Error("failed to get additional works", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get additional works")
}

func (a *WorkHandler) GetAdditionalWorkById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return errx.BadRequest(" id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		a.log.Error("failed to parse work id", zap.Error(err))
		return errx.BadRequest("work id must be a number")
	}

	res, err := a.service.GetAdditionalWorkById(id)
	if err != nil {
		a.log.Error("failed to get additional work by id", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get additional work by id")
}

func (a *WorkHandler) UpdateAdditionalWork(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return errx.BadRequest(" id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		a.log.Error("failed to parse work id", zap.Error(err))
		return errx.BadRequest("work id must be a number")
	}

	var request dto.UpdateAdditionalWorkRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	idCtx := c.Locals("accountId").(string)

	res, err := a.service.UpdateAdditionalWork(id, request, uuid.MustParse(idCtx))
	if err != nil {
		a.log.Error("failed to update additional work", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update additional work")
}

func (a *WorkHandler) DeleteAdditionalWork(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return errx.BadRequest(" id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		a.log.Error("failed to parse work id", zap.Error(err))
		return errx.BadRequest("work id must be a number")
	}

	err = a.service.DeleteAdditionalWork(id)
	if err != nil {
		a.log.Error("failed to delete additional work", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}

func (a *WorkHandler) GetOwnStaffWork(c *fiber.Ctx) error {
	idCtx := c.Locals("accountId").(string)

	res, err := a.service.GetStaffWorksByStaffId(uuid.MustParse(idCtx))
	if err != nil {
		a.log.Error("failed to get own staff work", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get own staff work")
}

func (a *WorkHandler) TakeAdditionalWork(c *fiber.Ctx) error {
	accountIdCtx := c.Locals("accountId").(string)

	idParam := c.Params("id")
	if idParam == "" {
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		a.log.Error("failed to parse work id", zap.Error(err))
		return errx.BadRequest("work id must be a number")
	}

	resp, err := a.service.TakeAdditionalWork(id, uuid.MustParse(accountIdCtx))
	if err != nil {
		a.log.Error("failed to take additional work", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, resp, "success take additional work")
}

func (a *WorkHandler) UpdateAdditionalWorkStaff(c *fiber.Ctx) error {
	accountIdCtx := c.Locals("accountId").(string)

	idParam := c.Params("id")
	if idParam == "" {
		return errx.BadRequest(" id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		a.log.Error("failed to parse work id", zap.Error(err))
		return errx.BadRequest("work id must be a number")
	}

	var request dto.UpdateAdditionalWorkStaffRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	res, err := a.service.UpdateAdditionalWorkStaff(id, request, uuid.MustParse(accountIdCtx))
	if err != nil {
		a.log.Error("failed to update additional work staff", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update additional work staff")
}

func (a *WorkHandler) UpdateDailyWorkStaff(c *fiber.Ctx) error {
	accountIdCtx := c.Locals("accountId").(string)

	IdParam := c.Params("id")
	if IdParam == "" {
		return errx.BadRequest(" id is required")
	}

	Id, err := strconv.ParseUint(IdParam, 10, 64)
	if err != nil {
		a.log.Error("failed to parse work id", zap.Error(err))
		return errx.BadRequest("work id must be a number")
	}

	var request dto.UpdateDailyWorkStaffRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	res, err := a.service.UpdateDailyWorkStaff(Id, request, uuid.MustParse(accountIdCtx))
	if err != nil {
		a.log.Error("failed to update daily work staff", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update daily work staff")
}

func (a *WorkHandler) DeleteDailyWork(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return errx.BadRequest(" id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		a.log.Error("failed to parse work id", zap.Error(err))
		return errx.BadRequest("work id must be a number")
	}

	err = a.service.DeleteDailyWork(id)
	if err != nil {
		a.log.Error("failed to delete daily work", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}
