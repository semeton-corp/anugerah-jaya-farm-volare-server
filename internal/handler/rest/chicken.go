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

type ChickenHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IChickenService
}

func (h *ChickenHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/chickens")
	v1.Get("/overview", middleware.Authentication(), h.GetChickenOverview)

	v1.Post("/monitorings", middleware.Authentication(), h.CreateChickenMonitoring)
	v1.Get("/monitorings", middleware.Authentication(), h.GetChickenMonitorings)
	v1.Put("/monitorings/:id", middleware.Authentication(), h.UpdateChickenMonitoring)
	v1.Get("/monitorings/:id", middleware.Authentication(), h.GetChickenMonitoringById)
	v1.Delete("/monitorings/:id", middleware.Authentication(), h.DeleteChickenMonitoring)

	v1.Post("/healths/items", middleware.Authentication(), h.CreateChickenHealthItem)
	v1.Get("/healths/items", middleware.Authentication(), h.GetChickenHealthItems)
	v1.Get("/healths/items/:id", middleware.Authentication(), h.GetChickenHealthItemById)
	v1.Put("/healths/items/:id", middleware.Authentication(), h.UpdateChickenHealthItem)
	v1.Delete("/healths/items/:id", middleware.Authentication(), h.DeleteChickenHealthItem)

	v1.Post("/healths/monitorings", middleware.Authentication(), h.CreateChickenHealthMonitoring)
	v1.Get("/healths/monitorings/details/:chickenCageId", middleware.Authentication(), h.GetChickenHealthMonitoringDetails)
	v1.Get("/healths/monitorings/:id", middleware.Authentication(), h.GetChickenHealthMonitoringById)
	v1.Put("/healths/monitorings/:id", middleware.Authentication(), h.UpdateChickenHealthMonitoring)
	v1.Delete("/healths/monitorings/:id", middleware.Authentication(), h.DeleteChickenHealthMonitoring)

	v1.Post("/procurements/drafts", middleware.Authentication(), h.CreateChickenProcurementDraft)
	v1.Put("/procurements/drafts/:id", middleware.Authentication(), h.UpdateChickenProcurementDraft)
	v1.Get("/procurements/drafts/:id", middleware.Authentication(), h.GetChickenProcurementDraft)
	v1.Delete("/procurements/drafts/:id", middleware.Authentication(), h.DeleteChickenProcurementDraft)
	v1.Get("/procurements/drafts", middleware.Authentication(), h.GetChickenProcurementDrafts)
	v1.Post("/procurements/drafts/:id/confirmations", middleware.Authentication(), h.ConfirmationChickenProcurementDraft)

	v1.Get("/procurements", middleware.Authentication(), h.GetChickenProcurements)
	v1.Get("/procurements/:id", middleware.Authentication(), h.GetChickenProcurement)
	v1.Post("/procurements/:id/arrivals", middleware.Authentication(), h.ArrivalConfirmationChickenProcurement)

	v1.Post("/procurements/:chickenProcurementId/payments", middleware.Authentication(), h.CreateChickenProcurementPayment)
	v1.Put("/procurements/:chickenProcurementId/payments/:id", middleware.Authentication(), h.UpdateChickenProcurementPayment)
	v1.Delete("/procurements/:chickenProcurementId/payments/:id", middleware.Authentication(), h.DeleteChickenProcurementPayment)

	v1.Post("/afkir/customers", middleware.Authentication(), h.CreateAfkirChickenCustomer)
	v1.Get("/afkir/customers", middleware.Authentication(), h.GetAfkirChickenCustomers)
	v1.Get("/afkir/customers/:id", middleware.Authentication(), h.GetAfkirChickenCustomer)
	v1.Put("/afkir/customers/:id", middleware.Authentication(), h.UpdateAfkirChickenCustomer)
	v1.Delete("/afkir/customers/:id", middleware.Authentication(), h.DeleteAfkirChickenCustomer)

	v1.Post("/afkir/sales/drafts", middleware.Authentication(), h.CreateAfkirChickenSaleDraft)
	v1.Get("/afkir/sales/drafts", middleware.Authentication(), h.GetAfkirChickenSaleDrafts)
	v1.Get("/afkir/sales/drafts/:id", middleware.Authentication(), h.GetAfkirChickenSaleDraft)
	v1.Put("/afkir/sales/drafts/:id", middleware.Authentication(), h.UpdateAfkirChickenSaleDraft)
	v1.Post("/afkir/sales/drafts/:id/confirmations", middleware.Authentication(), h.ConfirmationAfkirChickenSaleDraft)
	v1.Delete("/afkir/sales/drafts/:id", middleware.Authentication(), h.DeleteAfkirChickenSaleDraft)

	v1.Post("/afkir/sales", middleware.Authentication(), h.CreateAfkirChickenSale)
	v1.Get("/afkir/sales", middleware.Authentication(), h.GetAfkirChickenSales)
	v1.Get("/afkir/sales/:id", middleware.Authentication(), h.GetAfkirChickenSale)

	v1.Post("/afkir/sales/:afkirChickenSaleId/payments", middleware.Authentication(), h.CreateAfkirChickenSalePayment)
	v1.Put("/afkir/sales/:afkirChickenSaleId/payments/:id", middleware.Authentication(), h.UpdateAfkirChickenSalePayment)
	v1.Delete("/afkir/sales/:afkirChickenSaleId/payments/:id", middleware.Authentication(), h.DeleteAfkirChickenSalePayment)

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
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("userId not found in context")
		return errx.Unauthorized("userId not found in context")
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("validation failed", zap.Error(err))
		return err
	}

	res, err := h.service.CreateChickenMonitoring(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create chicken monitoring")
}

func (h *ChickenHandler) GetChickenMonitorings(c *fiber.Ctx) error {
	var filter dto.GetChickenMonitoringFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	res, err := h.service.GetChickenMonitorings(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get chicken monitorings")
}

func (h *ChickenHandler) GetChickenMonitoringById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	res, err := h.service.GetChickenMonitoringById(id)
	if err != nil {
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
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	var request dto.UpdateChickenMonitoringRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("validation failed", zap.Error(err))
		return err
	}

	res, err := h.service.UpdateChickenMonitoring(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		res,
		"success update chicken monitoring",
	)
}

func (h *ChickenHandler) DeleteChickenMonitoring(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	err = h.service.DeleteChickenMonitoring(id)
	if err != nil {
		h.log.Error("failed to delete chicken monitoring", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}

func (h *ChickenHandler) CreateChickenHealthItem(c *fiber.Ctx) error {
	var request dto.CreateChickenHealthItemRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validaton error", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.CreateChickenHealthItem(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, data, "success create chicken health item")
}

func (h *ChickenHandler) GetChickenHealthItemById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	data, err := h.service.GetChickenHealthItemById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success to get chicken health item by id")
}

func (h *ChickenHandler) GetChickenHealthItems(c *fiber.Ctx) error {
	var filter dto.GetChickenHealthItemFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	data, err := h.service.GetChickenHealthItems(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success to get chicken health items")
}

func (h *ChickenHandler) UpdateChickenHealthItem(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	var request dto.UpdateChickenHealthItemRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.UpdateChickenHealthItem(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success update chicken health item")
}

func (h *ChickenHandler) DeleteChickenHealthItem(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id param", zap.Error(err))
		return err
	}

	err = h.service.DeleteChickenHealthItem(id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *ChickenHandler) CreateChickenHealthMonitoring(c *fiber.Ctx) error {
	var request dto.CreateChickenHealthMonitoringRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Warn("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.CreateChickenHealthMonitoring(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, data, "success create chicken health monitoring")
}

func (h *ChickenHandler) GetChickenHealthMonitoringDetails(c *fiber.Ctx) error {
	chickenCageId, err := strconv.ParseUint(c.Params("chickenCageId"), 10, 64)
	if err != nil {
		h.log.Error("invalid chicken cage id param", zap.Error(err))
		return errx.BadRequest("invalid chicken cage id param")
	}

	data, err := h.service.GetChickenHealthMonitoringDetails(chickenCageId)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get chicken health monitoring details")
}

func (h *ChickenHandler) GetChickenHealthMonitoringById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	data, err := h.service.GetChickenHealthMonitoringById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get chicken health monitoring by id")
}

func (h *ChickenHandler) UpdateChickenHealthMonitoring(c *fiber.Ctx) error {
	var request dto.UpdateChickenHealthMonitoringRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation failed", zap.Error(err))
		return err
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Warn("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.UpdateChickenHealthMonitoring(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success chicken health monitoring")
}

func (h *ChickenHandler) DeleteChickenHealthMonitoring(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	err = h.service.DeleteChickenHealthMonitoring(id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *ChickenHandler) GetChickenOverview(c *fiber.Ctx) error {
	var filter dto.GetChickenOverviewFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query param", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	data, err := h.service.GetChickenOverview(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get chicken overview")
}

func (h *ChickenHandler) CreateChickenProcurementDraft(c *fiber.Ctx) error {
	var request dto.CreateChickenProcurementDraftRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.CreateChickenProcurementDraft(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, data, "success created chicken procurement draft")
}

func (h *ChickenHandler) GetChickenProcurementDrafts(c *fiber.Ctx) error {
	data, err := h.service.GetChickenProcurementDrafts()
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "get chicken procurement drafts")
}

func (h *ChickenHandler) GetChickenProcurementDraft(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed get chicken procurement draft", zap.Error(err))
		return err
	}

	data, err := h.service.GetChickenProcurementDraft(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "get chicken procurement draft")
}

func (h *ChickenHandler) UpdateChickenProcurementDraft(c *fiber.Ctx) error {
	var request dto.UpdateChickenProcurementDraftRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return err
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.UpdateChickenProcurementDraft(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success update chicken procurement draft")
}

func (h *ChickenHandler) DeleteChickenProcurementDraft(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("invaid id")
	}

	err = h.service.DeleteChickenProcurementDraft(id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *ChickenHandler) ConfirmationChickenProcurementDraft(c *fiber.Ctx) error {
	var request dto.ConfirmationChickenProcurementRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("invaid id")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found context")
		return errx.BadRequest("user id not found in context")
	}

	data, err := h.service.ConfirmationChickenProcurementDraft(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, data, "success confirm chicken procurement draft")
}

func (h *ChickenHandler) ArrivalConfirmationChickenProcurement(c *fiber.Ctx) error {
	var request dto.ArrivalConfirmationChickenProcurementRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("invaid id")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found context")
		return errx.BadRequest("user id not found in context")
	}

	data, err := h.service.ArrivalConfirmationChickenProcurement(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success arrival confirmation chicken procurement")
}

func (h *ChickenHandler) GetChickenProcurements(c *fiber.Ctx) error {
	var filter dto.GetChickenProcurementFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed parse filter", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&filter); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return err
	}

	data, err := h.service.GetChickenProcurements(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get chicken procurements")
}

func (h *ChickenHandler) GetChickenProcurement(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed parse request", zap.Error(err))
		return err
	}

	data, err := h.service.GetChickenProcurement(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get chicken procurement")
}

func (h *ChickenHandler) CreateChickenProcurementPayment(c *fiber.Ctx) error {
	var request dto.CreateChickenProcurementPaymentRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	chickenProcurementId, err := strconv.ParseUint(c.Params("chickenProcurementId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse chicken procurement id", zap.Error(err))
		return errx.BadRequest("invaid chicken procurement id")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found context")
		return errx.BadRequest("user id not found in context")
	}

	data, err := h.service.CreateChickenProcurementPayment(chickenProcurementId, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, data, "success create chicken procurement payment")
}

func (h *ChickenHandler) UpdateChickenProcurementPayment(c *fiber.Ctx) error {
	var request dto.UpdateChickenProcurementPaymentRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found context")
		return errx.BadRequest("user id not found in context")
	}

	chickenProcurementId, err := strconv.ParseUint(c.Params("chickenProcurementId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse chicken procurement id", zap.Error(err))
		return errx.BadRequest("invaid chicken procurement id")
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("invaid id")
	}

	data, err := h.service.UpdateChickenProcurementPayment(chickenProcurementId, id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, data, "success update chicken procurement payment")
}

func (h *ChickenHandler) DeleteChickenProcurementPayment(c *fiber.Ctx) error {
	chickenProcurementId, err := strconv.ParseUint(c.Params("chickenProcurementId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse chicken procurement id", zap.Error(err))
		return errx.BadRequest("invaid chicken procurement id")
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("invaid id")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found context")
		return errx.BadRequest("user id not found in context")
	}

	err = h.service.DeleteChickenProcurementPayment(chickenProcurementId, id, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *ChickenHandler) CreateAfkirChickenCustomer(c *fiber.Ctx) error {
	var request dto.CreateAfkirChickenCustomerRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found context")
		return errx.BadRequest("user id not found in context")
	}

	data, err := h.service.CreateAfkirChickenCustomer(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, data, "success create afkir chicken customer")
}

func (h *ChickenHandler) GetAfkirChickenCustomers(c *fiber.Ctx) error {
	data, err := h.service.GetAfkirChickenCustomers()
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get afkir chicken customers")
}

func (h *ChickenHandler) GetAfkirChickenCustomer(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("invaid id")
	}

	data, err := h.service.GetAfkirChickenCustomer(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success afkir chicken customer")
}

func (h *ChickenHandler) UpdateAfkirChickenCustomer(c *fiber.Ctx) error {
	var request dto.UpdateAfkirChickenCustomerRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("invaid id")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found context")
		return errx.BadRequest("user id not found in context")
	}

	data, err := h.service.UpdateAfkirChickenCustomer(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success create afkir chicken customer")
}

func (h *ChickenHandler) DeleteAfkirChickenCustomer(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("invaid id")
	}

	err = h.service.DeleteAfkirChickenCustomer(id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *ChickenHandler) CreateAfkirChickenSaleDraft(c *fiber.Ctx) error {
	var request dto.CreateAfkirChickenSaleDraftRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.CreateAkfirChickenSaleDraft(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, data, "success create afkir chicken sale draft")
}

func (h *ChickenHandler) GetAfkirChickenSaleDrafts(c *fiber.Ctx) error {
	data, err := h.service.GetAfkirChickenSaleDrafts()
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get chicken sale drafts")
}

func (h *ChickenHandler) GetAfkirChickenSaleDraft(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed parse id", zap.
			Error(err))
		return err
	}

	data, err := h.service.GetAfkirChickenSaleDraft(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get chicken sale draft")
}

func (h *ChickenHandler) UpdateAfkirChickenSaleDraft(c *fiber.Ctx) error {
	var request dto.UpdateAfkirChickenSaleDraftRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return errx.Unauthorized("user id not found in context")
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed parse id", zap.Error(err))
		return err
	}

	data, err := h.service.UpdateAfkirChickenSaleDraft(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success create afkir chicken sale draft")
}

func (h *ChickenHandler) DeleteAfkirChickenSaleDraft(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed parse id", zap.
			Error(err))
		return err
	}

	err = h.service.DeleteAfkirChickenSaleDraft(id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *ChickenHandler) CreateAfkirChickenSale(c *fiber.Ctx) error {
	var request dto.CreateAfkirChickenSaleRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.CreateAfkirChickenSale(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, data, "success create afkir chicken sale")
}

func (h *ChickenHandler) GetAfkirChickenSales(c *fiber.Ctx) error {
	var filter dto.GetAfkirChickenSaleFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed get afkir chicken sales", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&filter); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return err
	}

	data, err := h.service.GetAfkirChickenSales(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get afkir chicken sale")
}

func (h *ChickenHandler) GetAfkirChickenSale(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed parse id", zap.Error(err))
		return err
	}

	data, err := h.service.GetAkfirChickenSale(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get afkir chicken sale")
}

func (h *ChickenHandler) CreateAfkirChickenSalePayment(c *fiber.Ctx) error {
	var request dto.CreateAfkirChickenSalePaymentRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return err
	}

	afkirChickenSaleId, err := strconv.ParseUint(c.Params("afkirChickenSaleId"), 10, 64)
	if err != nil {
		h.log.Error("failed parse afkir chicken sale id", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.CreateAfkirChickenSalePayment(afkirChickenSaleId, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, data, "success create afkir chicken sale payment")
}

func (h *ChickenHandler) UpdateAfkirChickenSalePayment(c *fiber.Ctx) error {
	var request dto.UpdateAfkirChickenSalePaymentRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return err
	}

	afkirChickenSaleId, err := strconv.ParseUint(c.Params("afkirChickenSaleId"), 10, 64)
	if err != nil {
		h.log.Error("failed parse afkir chicken sale id", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return errx.Unauthorized("user id not found in context")
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed parse id", zap.Error(err))
		return err
	}

	data, err := h.service.UpdateAfkirChickenSalePayment(afkirChickenSaleId, id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success create afkir chicken sale payment")
}

func (h *ChickenHandler) DeleteAfkirChickenSalePayment(c *fiber.Ctx) error {
	afkirChickenSaleId, err := strconv.ParseUint(c.Params("afkirChickenSaleId"), 10, 64)
	if err != nil {
		h.log.Error("failed parse afkir chicken sale id", zap.Error(err))
		return err
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed parse id", zap.Error(err))
		return err
	}

	err = h.service.DeleteAfkirChickenSalePayment(afkirChickenSaleId, id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *ChickenHandler) ConfirmationAfkirChickenSaleDraft(c *fiber.Ctx) error {
	var request dto.CreateAfkirChickenSaleRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return errx.Unauthorized("user id not found in context")
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	data, err := h.service.ConfirmationAfkirChickenSaleDraft(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, data, "success confirmation afkir chicken sale")
}
