package rest

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/middleware"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/service"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/response"
	"go.uber.org/zap"
)

type ChickenHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IChickenService
}

func (a *ChickenHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/chickens")
	v1.Post("/monitorings", middleware.Authentication(constant.RoleAdmin), a.CreateChickenMonitoring)
	v1.Get("/monitorings", middleware.Authentication(constant.RoleAdmin), a.GetChickenMonitorings)
	v1.Put("/monitorings/:id", middleware.Authentication(constant.RoleAdmin), a.UpdateChickenMonitoring)
	v1.Get("/monitorings/:id", middleware.Authentication(constant.RoleAdmin), a.GetChickenMonitoringById)

	v1.Post("/monitorings/:chickenMonitoringId/diseases", middleware.Authentication(constant.RoleAdmin), a.CreateChickenDiseaseMonitoring)
	v1.Put("/monitorings/:chickenMonitoringId/diseases/:id", middleware.Authentication(constant.RoleAdmin), a.UpdateChickenDiseaseMonitoring)
	v1.Post("/monitorings/:chickenMonitoringId/vaccines", middleware.Authentication(constant.RoleAdmin), a.CreateChickenVacccineMonitoring)
	v1.Put("/monitorings/:chickenMonitoringId/vaccines/:id", middleware.Authentication(constant.RoleAdmin), a.UpdateChickenVaccineMonitoring)

	v1.Delete("/monitorings/:chickenMonitoringId/diseases/:id", middleware.Authentication(constant.RoleAdmin), a.DeleteChickenDiseaseMonitoring)
	v1.Delete("/monitorings/:chickenMonitoringId/vaccines/:id", middleware.Authentication(constant.RoleAdmin), a.DeleteChickenVaccineMonitoring)
	v1.Delete("/monitorings/:id", middleware.Authentication(constant.RoleAdmin), a.DeleteChickenMonitoring)

}

func NewChickenHandler(log *zap.Logger, service service.IChickenService, validator *validator.Validate) *ChickenHandler {
	return &ChickenHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *ChickenHandler) CreateChickenMonitoring(c *fiber.Ctx) error {
	var request dto.CreateChickenMonitoringRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[CreateChickenMonitoring] failed to parse request", zap.Error(err))
		return err
	}

	accountIdContext, ok := c.Locals("accountId").(string)
	if !ok {
		h.log.Error("[CreateChickenMonitoring] accountId not found in locals")
		return errx.Unauthorized("accountId not found in locals")
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[CreateChickenMonitoring] validation failed", zap.Error(err))
		return err
	}

	res, err := h.service.CreateChickenMonitoring(request, uuid.MustParse(accountIdContext))
	if err != nil {
		h.log.Error("[CreateChickenMonitoring] failed to create chicken monitoring", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		res,
		"success create chicken monitoring",
	)
}

func (h *ChickenHandler) GetChickenMonitorings(c *fiber.Ctx) error {
	var filter dto.GetChickenMonitoringFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("[GetChickens] failed to parse query", zap.Error(err))
		return err
	}

	res, err := h.service.GetChickenMonitorings(filter)
	if err != nil {
		h.log.Error("[GetChickens] failed to get chickens", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		res,
		"success get chickens",
	)
}

func (h *ChickenHandler) GetChickenMonitoringById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[GetChickenMonitoringById] id not found in params")
		return errx.BadRequest("id not found in params")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("[GetChickenMonitoringById] failed to parse id", zap.Error(err))
		return err
	}

	res, err := h.service.GetChickenMonitoringById(id)
	if err != nil {
		h.log.Error("[GetChickenMonitoringById] failed to get chicken monitoring", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		res,
		"success get chicken monitoring",
	)
}

func (h *ChickenHandler) UpdateChickenMonitoring(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[UpdateChickenMonitoring] id not found in params")
		return errx.BadRequest("id not found in params")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("[UpdateChickenMonitoring] failed to parse id", zap.Error(err))
		return err
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		h.log.Error("[ChangePassword] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	var request dto.UpdateChickenMonitoringRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[UpdateChickenMonitoring] failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[UpdateChickenMonitoring] validation failed", zap.Error(err))
		return err
	}

	res, err := h.service.UpdateChickenMonitoring(id, request, uuid.MustParse(accountId))
	if err != nil {
		h.log.Error("[UpdateChickenMonitoring] failed to update chicken monitoring", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		res,
		"success update chicken monitoring",
	)
}

func (h *ChickenHandler) CreateChickenDiseaseMonitoring(c *fiber.Ctx) error {
	var request dto.CreateChickenDiseaseMonitoringRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[CreateChickenDiseaseMonitoring] failed to parse request", zap.Error(err))
		return err
	}

	chickenMonitoringIdParam := c.Params("chickenMonitoringId")
	if chickenMonitoringIdParam == "" {
		h.log.Error("[CreateChickenDiseaseMonitoring] chickenMonitoringId not found in params")
		return errx.BadRequest("chickenMonitoringId not found in params")
	}

	chickenMonitoringId, err := strconv.ParseUint(chickenMonitoringIdParam, 10, 64)
	if err != nil {
		h.log.Error("[CreateChickenDiseaseMonitoring] failed to parse chickenMonitoringId", zap.Error(err))
		return err
	}

	accountIdContext, ok := c.Locals("accountId").(string)
	if !ok {
		h.log.Error("[CreateChickenDiseaseMonitoring] accountId not found in locals")
		return errx.Unauthorized("accountId not found in locals")
	}

	accountId, err := uuid.Parse(accountIdContext)
	if err != nil {
		h.log.Error("[CreateChickenDiseaseMonitoring] failed to parse accountId", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[CreateChickenDiseaseMonitoring] validation failed", zap.Error(err))
		return err
	}

	res, err := h.service.CreateChickenDiseaseMonitoring(chickenMonitoringId, request, accountId)
	if err != nil {
		h.log.Error("[CreateChickenDiseaseMonitoring] failed to create chicken disease monitoring", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusCreated,
		res,
		"success create chicken disease monitoring",
	)
}

func (h *ChickenHandler) CreateChickenVacccineMonitoring(c *fiber.Ctx) error {
	var request dto.CreateChickenVaccineMonitoringRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[CreateChickenVaccineMonitoring] failed to parse request", zap.Error(err))
		return err
	}

	chickenMonitoringIdParam := c.Params("chickenMonitoringId")
	if chickenMonitoringIdParam == "" {
		h.log.Error("[CreateChickenVaccineMonitoring] chickenMonitoringId not found in params")
		return errx.BadRequest("chickenMonitoringId not found in params")
	}

	chickenMonitoringId, err := strconv.ParseUint(chickenMonitoringIdParam, 10, 64)
	if err != nil {
		h.log.Error("[CreateChickenVaccineMonitoring] failed to parse chickenMonitoringId", zap.Error(err))
		return err
	}

	accountIdContext, ok := c.Locals("accountId").(string)
	if !ok {
		h.log.Error("[CreateChickenVaccineMonitoring] accountId not found in locals")
		return errx.Unauthorized("accountId not found in locals")
	}

	accountId, err := uuid.Parse(accountIdContext)
	if err != nil {
		h.log.Error("[CreateChickenVaccineMonitoring] failed to parse accountId", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[CreateChickenVaccineMonitoring] validation failed", zap.Error(err))
		return err
	}

	res, err := h.service.CreateChickenVaccineMonitoring(chickenMonitoringId, request, accountId)
	if err != nil {
		h.log.Error("[CreateChickenVaccineMonitoring] failed to create chicken vaccine monitoring", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusCreated,
		res,
		"success create chicken vaccine monitoring",
	)
}

func (h *ChickenHandler) UpdateChickenDiseaseMonitoring(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[UpdateChickenDiseaseMonitoring] id not found in params")
		return errx.BadRequest("id not found in params")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("[UpdateChickenDiseaseMonitoring] failed to parse id", zap.Error(err))
		return err
	}

	accountIdContext, ok := c.Locals("accountId").(string)
	if !ok {
		h.log.Error("[UpdateChickenDiseaseMonitoring] accountId not found in locals")
		return errx.Unauthorized("accountId not found in locals")
	}

	accountId, err := uuid.Parse(accountIdContext)
	if err != nil {
		h.log.Error("[UpdateChickenDiseaseMonitoring] failed to parse accountId", zap.Error(err))
		return err
	}

	var request dto.UpdateChickenDiseaseMonitoringRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[UpdateChickenDiseaseMonitoring] failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[UpdateChickenDiseaseMonitoring] validation failed", zap.Error(err))
		return err
	}

	res, err := h.service.UpdateChickenDiseaseMonitoring(id, request, accountId)
	if err != nil {
		h.log.Error("[UpdateChickenDiseaseMonitoring] failed to update chicken disease monitoring", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		res,
		"success update chicken disease monitoring",
	)
}

func (h *ChickenHandler) UpdateChickenVaccineMonitoring(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[UpdateChickenVaccineMonitoring] id not found in params")
		return errx.BadRequest("id not found in params")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("[UpdateChickenVaccineMonitoring] failed to parse id", zap.Error(err))
		return err
	}

	accountIdContext, ok := c.Locals("accountId").(string)
	if !ok {
		h.log.Error("[UpdateChickenVaccineMonitoring] accountId not found in locals")
		return errx.Unauthorized("accountId not found in locals")
	}

	accountId, err := uuid.Parse(accountIdContext)
	if err != nil {
		h.log.Error("[UpdateChickenVaccineMonitoring] failed to parse accountId", zap.Error(err))
		return err
	}

	var request dto.UpdateChickenVaccineMonitoringRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[UpdateChickenVaccineMonitoring] failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[UpdateChickenVaccineMonitoring] validation failed", zap.Error(err))
		return err
	}

	res, err := h.service.UpdateChickenVaccineMonitoring(id, request, accountId)
	if err != nil {
		h.log.Error("[UpdateChickenVaccineMonitoring] failed to update chicken vaccine monitoring", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		res,
		"success update chicken vaccine monitoring",
	)
}

func (h *ChickenHandler) DeleteChickenDiseaseMonitoring(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[DeleteChickenDiseaseMonitoring] id not found in params")
		return errx.BadRequest("id not found in params")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("[DeleteChickenDiseaseMonitoring] failed to parse id", zap.Error(err))
		return err
	}

	err = h.service.DeleteChickenDiseaseMonitoring(id)
	if err != nil {
		h.log.Error("[DeleteChickenDiseaseMonitoring] failed to delete chicken disease monitoring", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}

func (h *ChickenHandler) DeleteChickenVaccineMonitoring(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[DeleteChickenVaccineMonitoring] id not found in params")
		return errx.BadRequest("id not found in params")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("[DeleteChickenVaccineMonitoring] failed to parse id", zap.Error(err))
		return err
	}

	err = h.service.DeleteChickenVaccineMonitoring(id)
	if err != nil {
		h.log.Error("[DeleteChickenVaccineMonitoring] failed to delete chicken vaccine monitoring", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}

func (h *ChickenHandler) DeleteChickenMonitoring(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[DeleteChickenMonitoring] id not found in params")
		return errx.BadRequest("id not found in params")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("[DeleteChickenMonitoring] failed to parse id", zap.Error(err))
		return err
	}

	err = h.service.DeleteChickenMonitoring(id)
	if err != nil {
		h.log.Error("[DeleteChickenMonitoring] failed to delete chicken monitoring", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}
