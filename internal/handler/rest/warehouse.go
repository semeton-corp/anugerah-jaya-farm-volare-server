package rest

import (
	"net/http"
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

type WarehouseHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IWarehouseService
}

func (h *WarehouseHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/warehouses")

	v1.Post("/queues", middleware.Authentication(), h.CreateWarehouseSaleQueue)
	v1.Get("/queues", middleware.Authentication(), h.GetWarehouseSaleQueues)
	v1.Get("/queues/:id", middleware.Authentication(), h.GetWarehouseSaleQueue)
	v1.Delete("/queues/:id", middleware.Authentication(), h.DeleteWarehouseSaleQueue)
	v1.Post("/queues/:id/allocates", middleware.Authentication(), h.AllocateWarehouseSaleQueue)

	v1.Post("/sales", middleware.Authentication(), h.CreateWarehouseSale)
	v1.Get("/sales/:id", middleware.Authentication(), h.GetWarehouseSaleById)
	v1.Get("/sales", middleware.Authentication(), h.GetWarehouseSales)
	v1.Put("/sales/:id", middleware.Authentication(), h.UpdateWarehouseSale)
	v1.Delete("/sales/:id", middleware.Authentication(), h.DeleteWarehouseSale)
	v1.Post("/sales/:warehouseSaleId/payments", middleware.Authentication(), h.CreateWarehouseSalePayment)
	v1.Put("/sales/:warehouseSaleId/payments/:id", middleware.Authentication(), h.UpdateWarehouseSalePayment)
	v1.Delete("/sales/:warehouseSaleId/payments/:id", middleware.Authentication(), h.DeleteWarehouseSalePayment)
	v1.Patch("sales/:warehouseSaleId/send", middleware.Authentication(), h.SendWarehouseSale)

	v1.Post("/items", middleware.Authentication(), h.CreateWarehouseItem)
	v1.Get("/items", middleware.Authentication(), h.GetWarehouseItems)
	v1.Get("/:warehouseId/items/:itemId", middleware.Authentication(), h.GetWarehouseItemByWarehouseIdAndItemId)
	v1.Put("/:warehouseId/items/:itemId", middleware.Authentication(), h.UpdateWarehouseItem)
	v1.Delete("/:warehouseId/items/:itemId", middleware.Authentication(), h.DeleteWarehouseItem)

	v1.Get("/overview/:id", middleware.Authentication(), h.GetWarehouseOverview)
	v1.Post("/", middleware.Authentication(), h.CreateWarehouse)
	v1.Put("/:id", middleware.Authentication(), h.UpdateWarehouse)
	v1.Delete("/:id", middleware.Authentication(), h.DeleteWarehouse)
	v1.Get("/", middleware.Authentication(), h.GetWarehouses)
	v1.Get("/:id", middleware.Authentication(), h.GetWarehouseDetail)

	v1.Get("/items/eggs/summary/:warehouseId", middleware.Authentication(), h.GetEggWarehouseItemSummary)
	v1.Get("/items/corns/summary/:warehouseId", middleware.Authentication(), h.GetCornWarehouseItemSummary)

	v1.Get("/items/histories", middleware.Authentication(), h.GetWarehouseItemHistories)
	v1.Get("/items/histories/:id", middleware.Authentication(), h.GetWarehouseItemHistory)

	v1.Post("/items/procurements/drafts", middleware.Authentication(), h.CreateWarehouseItemProcurementDraft)
	v1.Get("/items/procurements/drafts/:id", middleware.Authentication(), h.GetWarehouseItemProcurementDraft)
	v1.Get("/items/procurements/drafts", middleware.Authentication(), h.GetWarehouseItemProcurementDrafts)
	v1.Put("/items/procurements/drafts/:id", middleware.Authentication(), h.UpdateWarehouseItemProcurementDraft)
	v1.Delete("/items/procurements/drafts/:id", middleware.Authentication(), h.DeleteWarehouseItemProcurementDraft)
	v1.Post("/items/procurements/drafts/:id/confirmations", middleware.Authentication(), h.ConfirmationWarehouseItemProcurementDraft)

	v1.Post("/items/procurements", middleware.Authentication(), h.CreateWarehouseItemProcurement)
	v1.Get("/items/procurements", middleware.Authentication(), h.GetWarehouseItemProcurements)
	v1.Get("/items/procurements/:id", middleware.Authentication(), h.GetWarehouseItemProcurement)
	v1.Put("/items/procurements/:id/arrivals", middleware.Authentication(), h.ArrivalConfirmationWarehouseItemProcurement)

	v1.Post("/items/procurements/:warehouseItemProcurementId/payments", middleware.Authentication(), h.CreateWarehouseItemProcurementPayment)
	v1.Put("/items/procurements/:warehouseItemProcurementId/payments/:id", middleware.Authentication(), h.UpdateWarehouseItemProcurementPayment)
	v1.Delete("/items/procurements/:warehouseItemProcurementId/payments/:id", middleware.Authentication(), h.DeleteWarehouseItemProcurementPayment)

	v1.Get("items/corns/prices", middleware.Authentication(), h.GetWarehouseItemCornPrices)

	v1.Post("/items/corns/procurements/drafts", middleware.Authentication(), h.CreateWarehouseItemCornProcurementDraft)
	v1.Get("/items/corns/procurements/drafts/:id", middleware.Authentication(), h.GetWarehouseItemCornProcurementDraft)
	v1.Get("/items/corns/procurements/drafts", middleware.Authentication(), h.GetWarehouseItemCornProcurementDrafts)
	v1.Put("/items/corns/procurements/drafts/:id", middleware.Authentication(), h.UpdateWarehouseItemCornProcurementDraft)
	v1.Delete("/items/corns/procurements/drafts/:id", middleware.Authentication(), h.DeleteWarehouseItemCornProcurementDraft)
	v1.Post("/items/corns/procurements/drafts/:id/confirmations", middleware.Authentication(), h.ConfirmationWarehouseItemCornProcurementDraft)

	v1.Post("/items/corns/procurements", middleware.Authentication(), h.CreateWarehouseItemCornProcurement)
	v1.Get("/items/corns/procurements", middleware.Authentication(), h.GetWarehouseItemCornProcurements)
	v1.Get("/items/corns/procurements/:id", middleware.Authentication(), h.GetWarehouseItemCornProcurement)
	v1.Put("/items/corns/procurements/:id/arrivals", middleware.Authentication(), h.ArrivalConfirmationWarehouseItemCornProcurement)

	v1.Post("/items/corns/procurements/:warehouseItemCornProcurementId/payments", middleware.Authentication(), h.CreateWarehouseItemCornProcurementPayment)
	v1.Put("/items/corns/procurements/:warehouseItemCornProcurementId/payments/:id", middleware.Authentication(), h.UpdateWarehouseItemCornProcurementPayment)
	v1.Delete("/items/corns/procurements/:warehouseItemCornProcurementId/payments/:id", middleware.Authentication(), h.DeleteWarehouseItemCornProcurementPayment)
}

func NewWarehouseHandler(log *zap.Logger, service service.IWarehouseService, validator *validator.Validate) *WarehouseHandler {
	return &WarehouseHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *WarehouseHandler) GetWarehouseOverview(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return err
	}

	data, err := h.service.GetWarehouseOverview(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get warehouse overview")
}

func (h *WarehouseHandler) GetWarehouseDetail(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return err
	}

	data, err := h.service.GetWarehouseDetailById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get warehouse detail")
}

func (h *WarehouseHandler) CreateWarehouse(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user if in context not found")
	}

	var request dto.CreateWarehouseRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("failed to validate struct", zap.Error(err))
		return err
	}

	res, err := h.service.CreateWarehouse(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, http.StatusCreated, res, "success to create warehouse")
}

func (h *WarehouseHandler) UpdateWarehouse(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse warehouse id")
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user if in context not found")
	}

	var request dto.UpdateWarehouseRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("failed to validate struct", zap.Error(err))
		return err
	}

	res, err := h.service.UpdateWarehouse(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, http.StatusCreated, res, "success to update warehouse")
}

func (h *WarehouseHandler) DeleteWarehouse(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse warehouse id")
		return err
	}

	err = h.service.DeleteWarehouse(id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *WarehouseHandler) GetWarehouses(c *fiber.Ctx) error {
	var filter dto.GetWarehouseFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query request", zap.Error(err))
		return err
	}

	warehouses, err := h.service.GetWarehouses(filter)
	if err != nil {
		h.log.Error("failed to get warehouses", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, warehouses, "get warehouses success")
}

func (h *WarehouseHandler) CreateWarehouseItem(c *fiber.Ctx) error {
	var request dto.CreateWarehouseItemRequest
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
		h.log.Error("failed to get user id from context")
		return errx.Unauthorized("no user id in context")
	}

	res, err := h.service.CreateWarehouseItem(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "create warehouse item success")
}

func (h *WarehouseHandler) GetWarehouseItems(c *fiber.Ctx) error {
	var filter dto.GetWarehouseItemFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	stockWarehouseItems, err := h.service.GetWarehouseItems(filter)
	if err != nil {
		h.log.Error("failed to get warehouse items", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, stockWarehouseItems, "get warehouse items success")
}

func (h *WarehouseHandler) GetWarehouseItemByWarehouseIdAndItemId(c *fiber.Ctx) error {
	warehouseIdParam := c.Params("warehouseId")
	warehouseItemIdParam := c.Params("itemId")

	if warehouseIdParam == "" || warehouseItemIdParam == "" {
		h.log.Error("warehouse id and item id are required")
		return errx.BadRequest("warrehouse id and item id are required")
	}

	warehouseId, err := strconv.ParseUint(warehouseIdParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse warehouse id", zap.Error(err))
		return errx.BadRequest("failed to parse warehouse id")
	}

	warehouseItemId, err := strconv.ParseUint(warehouseItemIdParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse item id", zap.Error(err))
		return errx.BadRequest("failed to parse warehouseItemId")
	}

	res, err := h.service.GetWarehouseItemByWarehouseIdAndItemId(warehouseId, warehouseItemId)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "get warehouse item success")
}

func (h *WarehouseHandler) UpdateWarehouseItem(c *fiber.Ctx) error {
	warehouseIdParam := c.Params("warehouseId")
	warehouseItemIdParam := c.Params("itemId")

	if warehouseIdParam == "" || warehouseItemIdParam == "" {
		h.log.Error("warehouseId and warehouseItemId are required")
		return errx.BadRequest("warehouseId and warehouseItemId are required")
	}

	warehouseId, err := strconv.ParseUint(warehouseIdParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse warehouseId", zap.Error(err))
		return errx.BadRequest("failed to parse warehouseId")
	}

	warehouseItemId, err := strconv.ParseUint(warehouseItemIdParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse warehouseItemId", zap.Error(err))
		return errx.BadRequest("failed to parse warehouseItemId")
	}

	var request dto.UpdateWarehouseItemRequest
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

	res, err := h.service.UpdateWarehouseItem(warehouseId, warehouseItemId, request, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("failed to update warehouse item", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "update warehouse item success")
}

func (h *WarehouseHandler) DeleteWarehouseItem(c *fiber.Ctx) error {
	warehouseIdParam := c.Params("warehouseId")
	warehouseItemIdParam := c.Params("itemId")

	if warehouseIdParam == "" || warehouseItemIdParam == "" {
		h.log.Error("warehouse id and warehouse item id are required")
		return errx.BadRequest("warehouse id and warehouse item id are required")
	}

	warehouseId, err := strconv.ParseUint(warehouseIdParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse warehouseId", zap.Error(err))
		return errx.BadRequest("failed to parse warehouse id")
	}

	warehouseItemId, err := strconv.ParseUint(warehouseItemIdParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse warehouseItemId", zap.Error(err))
		return errx.BadRequest("failed to parse warehouse item id")
	}

	err = h.service.DeleteWarehouseItem(warehouseId, warehouseItemId)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *WarehouseHandler) GetEggWarehouseItemSummary(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("warehouseId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse warehouse id")
		return errx.BadRequest("invalid warehouse id")
	}

	data, err := h.service.GetEggWarehouseItemSummary(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get egg warehouse item summary")
}

func (h *WarehouseHandler) GetCornWarehouseItemSummary(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("warehouseId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse warehouse id")
		return errx.BadRequest("invalid warehouse id")
	}

	data, err := h.service.GetCornWarehouseItemSummary(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get corn warehouse item summary")
}

func (h *WarehouseHandler) GetWarehouseItemHistories(c *fiber.Ctx) error {
	var filter dto.GetWarehouseItemHistoryFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query param", zap.Error(err))
		return err
	}

	data, err := h.service.GetWarehouseItemHistories(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get warehouse item histories")
}

func (h *WarehouseHandler) GetWarehouseItemHistory(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	data, err := h.service.GetWarehouseItemHistoryById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get warehouse item history")
}

func (h *WarehouseHandler) CreateWarehouseSale(c *fiber.Ctx) error {
	var request dto.CreateWarehouseSaleRequest
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

	res, err := h.service.CreateWarehouseSale(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create warehouse sale")
}

func (h *WarehouseHandler) GetWarehouseSaleById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	res, err := h.service.GetWarehouseSaleById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get warehouse sale by id")
}

func (h *WarehouseHandler) GetWarehouseSales(c *fiber.Ctx) error {
	var filter dto.GetWarehouseSaleFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	res, err := h.service.GetWarehouseSales(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get warehouse sales")
}

func (h *WarehouseHandler) CreateWarehouseSalePayment(c *fiber.Ctx) error {
	var request dto.CreateWarehouseSalePaymentRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	warehouseSaleId, err := strconv.ParseUint(c.Params("warehouseSaleId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse warehouse sale id", zap.Error(err))
		return errx.BadRequest("failed to parse warehouse sale id")
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("failed to get user id from context")
		return errx.Unauthorized("no user id in context")
	}

	res, err := h.service.CreateWarehouseSalePayment(warehouseSaleId, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create warehouse sale payment")
}

func (h *WarehouseHandler) UpdateWarehouseSale(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	var request dto.UpdateWarehouseSaleRequest
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

	res, err := h.service.UpdateWarehouseSale(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update warehouse sale")
}

func (h *WarehouseHandler) UpdateWarehouseSalePayment(c *fiber.Ctx) error {
	warehouseSaleId, err := strconv.ParseUint(c.Params("warehouseSaleId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse warehouse sale id", zap.Error(err))
		return errx.BadRequest("failed to parse warehouse sale id")
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	var request dto.UpdateWarehouseSalePaymentRequest
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

	res, err := h.service.UpdateWarehouseSalePayment(warehouseSaleId, id, request, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("failed to update warehouse sale payment", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update warehouse sale payment")
}

func (h *WarehouseHandler) DeleteWarehouseSalePayment(c *fiber.Ctx) error {
	storeSaleId, err := strconv.ParseUint(c.Params("warehouseSaleId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse  warehouse sale id", zap.Error(err))
		return errx.BadRequest("failed to warehouse sale id")
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("failed to get user id from context")
		return errx.Unauthorized("no user id in context")
	}

	err = h.service.DeleteWarehouseSalePayment(storeSaleId, id, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *WarehouseHandler) SendWarehouseSale(c *fiber.Ctx) error {
	warehouseSaleId, err := strconv.ParseUint(c.Params("warehouseSaleId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse warehouse sale id", zap.Error(err))
		return errx.BadRequest("failed to parse warehouse sale id")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	res, err := h.service.SendWarehouseSale(warehouseSaleId, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("failed to send warehouse sale", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success send warehouse sale")
}

func (h *WarehouseHandler) DeleteWarehouseSale(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("failed to get user id from context")
		return errx.Unauthorized("no user id in context")
	}

	err = h.service.DeleteWarehouseSale(id, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *WarehouseHandler) CreateWarehouseSaleQueue(c *fiber.Ctx) error {
	var request dto.CreateWarehouseSaleQueueRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed parse request bory", zap.Error(err))
		return errx.BadRequest(err.Error())
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.CreateWarehouseSaleQueue(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, data, "success created warehouse sale queue")
}

func (h *WarehouseHandler) GetWarehouseSaleQueues(c *fiber.Ctx) error {
	var filter dto.GetWarehouseSaleQueueFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failde parse query filter", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&filter); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return err
	}

	data, err := h.service.GetWarehouseSaleQueues(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get warehouse sale queues")
}

func (h *WarehouseHandler) GetWarehouseSaleQueue(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return err
	}

	data, err := h.service.GetWarehouseSaleQueue(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get store sale")
}

func (h *WarehouseHandler) DeleteWarehouseSaleQueue(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return err
	}

	err = h.service.DeleteWarehouseSaleQueue(id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *WarehouseHandler) AllocateWarehouseSaleQueue(c *fiber.Ctx) error {
	var request dto.CreateWarehouseSaleRequest
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
		h.log.Error("failed to parse id", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.AllocateWarehouseSaleQueue(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success allocate store sale queue")
}

func (h *WarehouseHandler) CreateWarehouseItemProcurementDraft(c *fiber.Ctx) error {
	var request dto.CreateWarehouseItemProcurementDraftRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return errx.BadRequest(err.Error())
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return errx.Unauthorized("user id not found in context")
	}

	res, err := h.service.CreateWarehouseItemProcurementDraft(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create warehouse item procurement draft")
}

func (h *WarehouseHandler) GetWarehouseItemProcurementDrafts(c *fiber.Ctx) error {
	res, err := h.service.GetWarehouseItemProcurementDrafts()
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusOK, res, "success get warehouse item procurement drafts")
}

func (h *WarehouseHandler) GetWarehouseItemProcurementDraft(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	res, err := h.service.GetWarehouseItemProcurementDraft(id)
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusOK, res, "success get warehouse item procurement draft")
}

func (h *WarehouseHandler) UpdateWarehouseItemProcurementDraft(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	var request dto.UpdateWarehouseItemProcurementDraftRequest
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

	res, err := h.service.UpdateWarehouseItemProcurementDraft(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusOK, res, "success update warehouse item procurement draft")
}

func (h *WarehouseHandler) ConfirmationWarehouseItemProcurementDraft(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	var request dto.CreateWarehouseItemProcurementRequest
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

	res, err := h.service.ConfirmationWarehouseItemProcurementDraft(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusCreated, res, "success confirmation warehouse item procurement draft")
}

func (h *WarehouseHandler) CreateWarehouseItemProcurement(c *fiber.Ctx) error {
	var request dto.CreateWarehouseItemProcurementRequest
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

	res, err := h.service.CreateWarehouseItemProcurement(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create warehouse item procurement")
}

func (h *WarehouseHandler) GetWarehouseItemProcurements(c *fiber.Ctx) error {
	var filter dto.GetWarehouseItemProcurementFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	res, err := h.service.GetWarehouseItemProcurements(filter)
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusOK, res, "success get warehouse item procurements")
}

func (h *WarehouseHandler) GetWarehouseItemProcurement(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	res, err := h.service.GetWarehouseItemProcurement(id)
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusOK, res, "success get warehouse item procurement")
}

func (h *WarehouseHandler) DeleteWarehouseItemProcurementDraft(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	err = h.service.DeleteWarehouseItemProcurementDraft(id)
	if err != nil {
		return err
	}
	return response.NoContentResponse(c)
}

func (h *WarehouseHandler) ArrivalConfirmationWarehouseItemProcurement(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	var request dto.ArrivalConfirmationWarehouseItemProcurementRequest
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

	res, err := h.service.ArrivalConfirmationWarehouseItemProcurement(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusOK, res, "success arrival confirmation warehouse item procurement")
}

func (h *WarehouseHandler) CreateWarehouseItemProcurementPayment(c *fiber.Ctx) error {
	warehouseItemProcurementId, err := strconv.ParseUint(c.Params("warehouseItemProcurementId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse warehouse item procurement id", zap.Error(err))
		return errx.BadRequest("failed to parse warehouse item procurement id")
	}

	var request dto.CreateWarehouseItemProcurementPaymentRequest
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

	res, err := h.service.CreateWarehouseItemProcurementPayment(warehouseItemProcurementId, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create warehouse item procurement payment")
}

func (h *WarehouseHandler) UpdateWarehouseItemProcurementPayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	warehouseItemProcurementId, err := strconv.ParseUint(c.Params("warehouseItemProcurementId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse warehouse item procurement id", zap.Error(err))
		return errx.BadRequest("failed to parse warehouse item procurement id")
	}

	var request dto.UpdateWarehouseItemProcurementPaymentRequest
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

	res, err := h.service.UpdateWarehouseItemProcurementPayment(id, warehouseItemProcurementId, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusOK, res, "success update warehouse item procurement payment")
}

func (h *WarehouseHandler) DeleteWarehouseItemProcurementPayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	warehouseItemProcurementId, err := strconv.ParseUint(c.Params("warehouseItemProcurementId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse warehouse item procurement id", zap.Error(err))
		return errx.BadRequest("failed to parse warehouse item procurement id")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return errx.Unauthorized("user id not found in context")
	}

	err = h.service.DeleteWarehouseItemProcurementPayment(id, warehouseItemProcurementId, uuid.MustParse(userId))
	if err != nil {
		return err
	}
	return response.NoContentResponse(c)
}

func (h *WarehouseHandler) CreateWarehouseItemCornProcurementDraft(c *fiber.Ctx) error {
	var request dto.CreateWarehouseItemCornProcurementDraftRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return errx.BadRequest(err.Error())
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return errx.Unauthorized("user id not found in context")
	}

	res, err := h.service.CreateWarehouseItemCornProcurementDraft(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create warehouse item corn procurement draft")
}

func (h *WarehouseHandler) GetWarehouseItemCornProcurementDrafts(c *fiber.Ctx) error {
	res, err := h.service.GetWarehouseItemCornProcurementDrafts()
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusOK, res, "success get warehouse item corn procurement drafts")
}

func (h *WarehouseHandler) GetWarehouseItemCornProcurementDraft(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	res, err := h.service.GetWarehouseItemCornProcurementDraft(id)
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusOK, res, "success get warehouse item corn procurement draft")
}

func (h *WarehouseHandler) UpdateWarehouseItemCornProcurementDraft(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	var request dto.UpdateWarehouseItemCornProcurementDraftRequest
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

	res, err := h.service.UpdateWarehouseItemCornProcurementDraft(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusOK, res, "success update warehouse item corn procurement draft")
}

func (h *WarehouseHandler) DeleteWarehouseItemCornProcurementDraft(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	err = h.service.DeleteWarehouseItemCornProcurementDraft(id)
	if err != nil {
		return err
	}
	return response.NoContentResponse(c)
}

func (h *WarehouseHandler) ConfirmationWarehouseItemCornProcurementDraft(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	var request dto.CreateWarehouseItemCornProcurementRequest
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

	res, err := h.service.ConfirmationWarehouseItemCornProcurementDraft(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusCreated, res, "success allocate warehouse item corn procurement draft")
}

func (h *WarehouseHandler) CreateWarehouseItemCornProcurement(c *fiber.Ctx) error {
	var request dto.CreateWarehouseItemCornProcurementRequest
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

	res, err := h.service.CreateWarehouseItemCornProcurement(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create warehouse item corn procurement")
}

func (h *WarehouseHandler) GetWarehouseItemCornProcurements(c *fiber.Ctx) error {
	var filter dto.GetWarehouseItemCornProcurementFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	res, err := h.service.GetWarehouseItemCornProcurements(filter)
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusOK, res, "success get warehouse item corn procurements")
}

func (h *WarehouseHandler) GetWarehouseItemCornProcurement(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	res, err := h.service.GetWarehouseItemCornProcurement(id)
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusOK, res, "success get warehouse item corn procurement")
}

func (h *WarehouseHandler) ArrivalConfirmationWarehouseItemCornProcurement(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	var request dto.ArrivalConfirmationWarehouseItemCornProcurementRequest
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

	res, err := h.service.ArrivalConfirmationWarehouseItemCornProcurement(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusOK, res, "success arrival confirmation warehouse item corn procurement")
}

func (h *WarehouseHandler) CreateWarehouseItemCornProcurementPayment(c *fiber.Ctx) error {
	warehouseItemCornProcurementId, err := strconv.ParseUint(c.Params("warehouseItemCornProcurementId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse warehouse item corn procurement id", zap.Error(err))
		return errx.BadRequest("failed to parse warehouse item corn procurement id")
	}

	var request dto.CreateWarehouseItemCornProcurementPaymentRequest
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

	res, err := h.service.CreateWarehouseItemCornProcurementPayment(warehouseItemCornProcurementId, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create warehouse item corn procurement payment")
}

func (h *WarehouseHandler) UpdateWarehouseItemCornProcurementPayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	warehouseItemCornProcurementId, err := strconv.ParseUint(c.Params("warehouseItemCornProcurementId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse warehouse item corn procurement id", zap.Error(err))
		return errx.BadRequest("failed to parse warehouse item corn procurement id")
	}

	var request dto.UpdateWarehouseItemCornProcurementPaymentRequest
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

	res, err := h.service.UpdateWarehouseItemCornProcurementPayment(id, warehouseItemCornProcurementId, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}
	return response.SuccessResponse(c, fiber.StatusOK, res, "success update warehouse item corn procurement payment")
}

func (h *WarehouseHandler) DeleteWarehouseItemCornProcurementPayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	warehouseItemCornProcurementId, err := strconv.ParseUint(c.Params("warehouseItemCornProcurementId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse warehouse item corn procurement id", zap.Error(err))
		return errx.BadRequest("failed to parse warehouse item corn procurement id")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return errx.Unauthorized("user id not found in context")
	}

	err = h.service.DeleteWarehouseItemCornProcurementPayment(id, warehouseItemCornProcurementId, uuid.MustParse(userId))
	if err != nil {
		return err
	}
	return response.NoContentResponse(c)
}

func (h *WarehouseHandler) GetWarehouseItemCornPrices(c *fiber.Ctx) error {
	data, err := h.service.GetWarehouseItemCornPrices()
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get warehouse item corn prices")
}
