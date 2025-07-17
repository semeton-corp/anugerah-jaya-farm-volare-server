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
	v1.Get("/overview/:id", middleware.Authentication(), h.GetWarehouseOverview)
	v1.Post("/", middleware.Authentication(), h.CreateWarehouse)
	v1.Put("/:id", middleware.Authentication(), h.UpdateWarehouse)
	v1.Delete("/:id", middleware.Authentication(), h.DeleteWarehouse)
	v1.Get("/", middleware.Authentication(), h.GetWarehouses)
	v1.Get("/:id", middleware.Authentication(), h.GetWarehouseDetail)

	v1.Get("/items/eggs/summary/:warehouseId", middleware.Authentication(), h.GetEggWarehouseItemSummary)
	v1.Post("/items", middleware.Authentication(), h.CreateWarehouseItem)
	v1.Get("/items", middleware.Authentication(), h.GetWarehouseItems)
	v1.Get("/:warehouseId/items/:itemId", middleware.Authentication(), h.GetWarehouseItemByWarehouseIdAndItemId)
	v1.Put("/:warehouseId/items/:itemId", middleware.Authentication(), h.UpdateWarehouseItem)
	v1.Delete("/:warehouseId/items/:itemId", middleware.Authentication(), h.DeleteWarehouseItem)

	v1.Post("/order/items", middleware.Authentication(), h.CreateWarehouseOrderItem)
	v1.Get("/order/items", middleware.Authentication(), h.GetWarehouseOrderItems)
	v1.Get("/order/items/:id", middleware.Authentication(), h.GetWarehouseOrderItemById)
	v1.Delete("/order/items/:id", middleware.Authentication(), h.DeleteWarehouseOrderItem)
	v1.Patch("/order/items/:id/takes", middleware.Authentication(), h.TakeWarehouseOrderItem)

	v1.Post("/sales", middleware.Authentication(), h.CreateWarehouseSale)
	v1.Get("/sales/:id", middleware.Authentication(), h.GetWarehouseSaleById)
	v1.Get("/sales", middleware.Authentication(), h.GetWarehouseSales)
	v1.Put("/sales/:id", middleware.Authentication(), h.UpdateWarehouseSale)
	v1.Post("/sales/:warehouseSaleId/payments", middleware.Authentication(), h.CreateWarehouseSalePayment)
	v1.Put("/sales/:warehouseSaleId/payments/:id", middleware.Authentication(), h.UpdateWarehouseSalePayment)
	v1.Patch("sales/:warehouseSaleId/send", middleware.Authentication(), h.SendWarehouseSale)

	v1.Get("/items/histories", middleware.Authentication(), h.GetWarehouseItemHistories)
	v1.Get("/items/histories/:id", middleware.Authentication(), h.GetWarehouseItemHistory)
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

func (h *WarehouseHandler) CreateWarehouseOrderItem(c *fiber.Ctx) error {
	var request dto.CreateWarehouseOrderItemRequest
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

	res, err := h.service.CreateWarehouseOrderItem(request, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("failed to create store order item", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "create store order item success")
}

func (h *WarehouseHandler) GetWarehouseOrderItemById(c *fiber.Ctx) error {
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

	res, err := h.service.GetWarehouseOrderItemById(id)
	if err != nil {
		h.log.Error("failed to get store order item", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "get store order item success")
}

func (h *WarehouseHandler) GetWarehouseOrderItems(c *fiber.Ctx) error {
	var filter dto.GetWarehouseOrderItemFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	warehouseOrderItems, err := h.service.GetWarehouseOrderItems(filter)
	if err != nil {
		h.log.Error("failed to get store order items", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, warehouseOrderItems, "get store order items success")
}

func (h *WarehouseHandler) DeleteWarehouseOrderItem(c *fiber.Ctx) error {
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

	err = h.service.DeleteWarehouseOrderItem(id)
	if err != nil {
		h.log.Error("failed to delete store order item", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}

func (h *WarehouseHandler) TakeWarehouseOrderItem(c *fiber.Ctx) error {
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

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	res, err := h.service.TakeWarehouseOrderItem(id, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("failed to take store order item", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "take store order item success")
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

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create store sale")
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

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get store sale by id")
}

func (h *WarehouseHandler) GetWarehouseSales(c *fiber.Ctx) error {
	var filter dto.GetWarehouseSaleFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("[GetStoreSales] failed to parse query", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("[GetStoreSales] failed to validate request", zap.Error(err))
		return err
	}

	res, err := h.service.GetWarehouseSales(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get store sales")
}

func (h *WarehouseHandler) CreateWarehouseSalePayment(c *fiber.Ctx) error {
	var request dto.CreateWarehouseSalePaymentRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[CreateStoreSalePayment] failed to parse request", zap.Error(err))
		return err
	}

	storeSaleIdParam := c.Params("storeSaleId")
	if storeSaleIdParam == "" {
		h.log.Error("storeSaleId is required")
		return errx.BadRequest("storeSaleId is required")
	}

	storeSaleId, err := strconv.ParseUint(storeSaleIdParam, 10, 64)
	if err != nil {
		h.log.Error("failed to parse storeSaleId", zap.Error(err))
		return errx.BadRequest("failed to parse storeSaleId")
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

	res, err := h.service.CreateWarehouseSalePayment(storeSaleId, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create store sale payment")
}

func (h *WarehouseHandler) UpdateWarehouseSale(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[UpdateStoreSale] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("[UpdateStoreSale] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	var request dto.UpdateWarehouseSaleRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[UpdateStoreSale] failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[UpdateStoreSale] failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[UpdateStoreSale] failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	res, err := h.service.UpdateWarehouseSale(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update store sale")
}

func (h *WarehouseHandler) UpdateWarehouseSalePayment(c *fiber.Ctx) error {
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

	res, err := h.service.UpdateWarehouseSalePayment(id, request, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[UpdateStoreSalePayment] failed to update store sale payment", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update store sale payment")
}

func (h *WarehouseHandler) SendWarehouseSale(c *fiber.Ctx) error {
	storeSaleIdParam := c.Params("storeSaleId")
	if storeSaleIdParam == "" {
		h.log.Error("[SendStoreSale] storeSaleId is required")
		return errx.BadRequest("storeSaleId is required")
	}

	storeSaleId, err := strconv.ParseUint(storeSaleIdParam, 10, 64)
	if err != nil {
		h.log.Error("[SendStoreSale] failed to parse storeSaleId", zap.Error(err))
		return errx.BadRequest("failed to parse storeSaleId")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[SendStoreSale] failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	res, err := h.service.SendWarehouseSale(storeSaleId, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[SendStoreSale] failed to send store sale", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success send store sale")
}
