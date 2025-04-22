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

type EggHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IEggService
}

func (a *EggHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/eggs")
	v1.Post("/monitorings", middleware.Authentication(), a.CreateEggMonitoring)
	v1.Get("/monitorings", middleware.Authentication(), a.GetEggMonitorings)
	v1.Get("/monitorings/:id", middleware.Authentication(), a.GetEggMonitoringById)
	v1.Put("/monitorings/:id", middleware.Authentication(), a.UpdateEggMonitoring)
	v1.Delete("/monitorings/:id", middleware.Authentication(), a.DeleteEggMonitoring)
}

func NewEggHandler(log *zap.Logger, service service.IEggService, validator *validator.Validate) *EggHandler {
	return &EggHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (a *EggHandler) CreateEggMonitoring(c *fiber.Ctx) error {
	var request dto.CreateEggMonitoringRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[CreateEggMonitoring] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[CreateEggMonitoring] failed to validate request", zap.Error(err))
		return err
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[CreateEggMonitoring] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	res, err := a.service.CreateEggMonitoring(request, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[CreateEggMonitoring] failed to create egg monitoring", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create egg monitoring")
}

func (a *EggHandler) GetEggMonitorings(c *fiber.Ctx) error {
	var filter dto.GetEggMonitoringFilter
	if err := c.QueryParser(&filter); err != nil {
		a.log.Error("[GetEggMonitorings] failed to parse query", zap.Error(err))
		return err
	}

	res, err := a.service.GetEggMonitorings(filter)
	if err != nil {
		a.log.Error("[GetEggMonitorings] failed to get egg monitorings", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get egg monitorings")
}

func (a *EggHandler) GetEggMonitoringById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		a.log.Error("[GetEggMonitoringById] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		a.log.Error("[GetEggMonitoringById] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	res, err := a.service.GetEggMonitoringById(id)
	if err != nil {
		a.log.Error("[GetEggMonitoringById] failed to get egg monitoring by id", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get egg monitoring by id")
}

func (a *EggHandler) UpdateEggMonitoring(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		a.log.Error("[UpdateEggMonitoring] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		a.log.Error("[UpdateEggMonitoring] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	var request dto.UpdateEggMonitoringRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[UpdateEggMonitoring] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[UpdateEggMonitoring] failed to validate request", zap.Error(err))
		return err
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[UpdateEggMonitoring] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	res, err := a.service.UpdateEggMonitoring(id, request, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[UpdateEggMonitoring] failed to update egg monitoring", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update egg monitoring")
}

func (a *EggHandler) DeleteEggMonitoring(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		a.log.Error("[DeleteEggMonitoring] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		a.log.Error("[DeleteEggMonitoring] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	if err := a.service.DeleteEggMonitoring(id); err != nil {
		a.log.Error("[DeleteEggMonitoring] failed to delete egg monitoring", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}
