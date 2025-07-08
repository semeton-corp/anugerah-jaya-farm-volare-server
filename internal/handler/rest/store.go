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

type StoreHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IStoreService
}

func (h *StoreHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/stores")
	v1.Get("/", middleware.Authentication(), h.GetStores)
	v1.Post("/", middleware.Authentication(), h.CreateStore)
	v1.Put("/:id", middleware.Authentication(), h.UpdateStore)
	v1.Delete("/:id", middleware.Authentication(), h.DeleteStore)
	v1.Get("/:id", middleware.Authentication(), h.GetStoreDetail)
	v1.Get("/overview/:id", middleware.Authentication(), h.GetStoreverview)

	v1.Post("/request/items", middleware.Authentication(), h.CreateStoreRequestItem)
	v1.Get("/request/items", middleware.Authentication(), h.GetStockRequestItems)
	v1.Get("/request/items/:id", middleware.Authentication(), h.GetStoreRequestItemById)
	v1.Put("/request/items/:id/warehouses", middleware.Authentication(), h.UpdateStoreRequestItemByWarehouse)
	v1.Put("/request/items/:id/stores", middleware.Authentication(), h.UpdateStoreRequestItemByStore)

	v1.Get("/items/overview/:id", middleware.Authentication(), h.GetStoreverview)
	v1.Get("/:storeId/items/:itemId", middleware.Authentication(), h.GetStoreItem)
	v1.Put("/:storeId/items/:itemId", middleware.Authentication(), h.UpdateStoreItem)

	v1.Post("/sales", middleware.Authentication(), h.CreateStoreSale)
	v1.Get("/sales/:id", middleware.Authentication(), h.GetStoreSaleById)
	v1.Get("/sales", middleware.Authentication(), h.GetStoreSales)
	v1.Put("/sales/:id", middleware.Authentication(), h.UpdateStoreSale)
	v1.Post("/sales/:storeSaleId/payments", middleware.Authentication(), h.CreateStoreSalePayment)
	v1.Put("/sales/:storeSaleId/payments/:id", middleware.Authentication(), h.UpdateStoreSalePayment)
	v1.Patch("sales/:storeSaleId/send", middleware.Authentication(), h.SendStoreSale)
}

func NewStoreHandler(log *zap.Logger, service service.IStoreService, validator *validator.Validate) *StoreHandler {
	return &StoreHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *StoreHandler) GetStoreDetail(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return err
	}

	data, err := h.service.GetStoreDetailById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get store detail")
}

func (h *StoreHandler) CreateStore(c *fiber.Ctx) error {
	var request dto.CreateStoreRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("validation failed", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.CreateStore(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, data, "success create store")
}

func (h *StoreHandler) UpdateStore(c *fiber.Ctx) error {
	var request dto.UpdateStoreRequest
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
		return err

	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.UpdateStore(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success update store")
}

func (h *StoreHandler) DeleteStore(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id param", zap.Error(err))
		return err
	}

	err = h.service.DeleteStore(id)
	if err != nil {
		h.log.Error("failed to delete store")
		return err
	}

	return response.NoContentResponse(c)
}

func (h *StoreHandler) GetStores(c *fiber.Ctx) error {
	var filter dto.GetStoreFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query request")
		return err
	}

	stores, err := h.service.GetStores(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, stores, "success get stores")
}

func (h *StoreHandler) CreateStoreRequestItem(c *fiber.Ctx) error {
	var request dto.CreateStoreRequestItemRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("failed to get user id from context")
		return errx.Unauthorized("no user id in context")
	}

	res, err := h.service.CreateStoreRequestItem(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create store request item")
}

func (h *StoreHandler) GetStockRequestItems(c *fiber.Ctx) error {
	var filter dto.GetStoreRequestItemFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("[GetStockRequestItems] failed to parse query", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("[GetStockRequestItems] failed to validate request", zap.Error(err))
		return err
	}

	res, err := h.service.GetStoreRequestItems(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get stock request items")
}

func (h *StoreHandler) GetStoreRequestItemById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[GetStoreRequestItemById] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("[GetStoreRequestItemById] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	res, err := h.service.GetStoreRequestItemById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get store request item by id")
}

func (h *StoreHandler) UpdateStoreRequestItemByWarehouse(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[UpdateStoreRequestItemByWarehouse] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("[UpdateStoreRequestItemByWarehouse] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	var request dto.UpdateStoreRequestItemByWarehouseRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[UpdateStoreRequestItemByWarehouse] failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[UpdateStoreRequestItemByWarehouse] failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[UpdateStoreRequestItemByWarehouse] failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	storeRequestItem, err := h.service.GetStoreRequestItemById(id)
	if err != nil {
		return err
	}

	combinedRequest := dto.UpdateStoreRequestItemRequest{
		Quantity: storeRequestItem.Quantity,
		Status:   request.Status,
	}

	res, err := h.service.UpdateStoreRequestItem(id, combinedRequest, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[UpdateStoreRequestItemByWarehouse] failed to update store request item by warehouse", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update store request item by warehouse")
}

func (h *StoreHandler) UpdateStoreRequestItemByStore(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[UpdateStoreRequestItemByWarehouse] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("[UpdateStoreRequestItemByWarehouse] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	var request dto.UpdateStoreRequestItemByStoreRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[UpdateStoreRequestItemByWarehouse] failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[UpdateStoreRequestItemByWarehouse] failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[UpdateStoreRequestItemByWarehouse] failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	combinedRequest := dto.UpdateStoreRequestItemRequest{
		Quantity: request.Quantity,
		Status:   request.Status,
	}

	res, err := h.service.UpdateStoreRequestItem(id, combinedRequest, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update store request item by warehouse")
}

func (h *StoreHandler) GetStoreverview(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	data, err := h.service.GetStoreOverview(id)
	if err != nil {
		h.log.Error("failed to get store item overview", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get store item overview")
}

func (h *StoreHandler) GetStoreItem(c *fiber.Ctx) error {
	storeId, err := strconv.ParseUint(c.Params("storeId"), 10, 64)
	if err != nil {
		h.log.Error("invalid store id param", zap.Error(err))
		return errx.BadRequest("invalid store id param")
	}

	itemId, err := strconv.ParseUint(c.Params("itemId"), 10, 64)
	if err != nil {
		h.log.Error("invalid item id param", zap.Error(err))
		return errx.BadRequest("invalid item id param")
	}

	data, err := h.service.GetStoreItemByStoreIdAndItemId(storeId, itemId)
	if err != nil {
		h.log.Error("failed to get store item overview", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get store item")
}

func (h *StoreHandler) UpdateStoreItem(c *fiber.Ctx) error {
	var request dto.UpdateStoreItemRequest

	storeId, err := strconv.ParseUint(c.Params("storeId"), 10, 64)
	if err != nil {
		h.log.Error("invalid store id param", zap.Error(err))
		return errx.BadRequest("invalid store id param")
	}

	itemId, err := strconv.ParseUint(c.Params("itemId"), 10, 64)
	if err != nil {
		h.log.Error("invalid item id param", zap.Error(err))
		return errx.BadRequest("invalid item id param")
	}

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
		h.log.Warn("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	data, err := h.service.UpdateStoreItem(storeId, itemId, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success update store item")
}

func (h *StoreHandler) CreateStoreSale(c *fiber.Ctx) error {
	var request dto.CreateStoreSaleRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[CreateStoreSale] failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[CreateStoreSale] failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[CreateStoreSale] failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	res, err := h.service.CreateStoreSale(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create store sale")
}

func (h *StoreHandler) GetStoreSaleById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[GetStoreSaleById] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("[GetStoreSaleById] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	res, err := h.service.GetStoreSaleById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get store sale by id")
}

func (h *StoreHandler) GetStoreSales(c *fiber.Ctx) error {
	var filter dto.GetStoreSaleFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("[GetStoreSales] failed to parse query", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("[GetStoreSales] failed to validate request", zap.Error(err))
		return err
	}

	res, err := h.service.GetStoreSales(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get store sales")
}

func (h *StoreHandler) CreateStoreSalePayment(c *fiber.Ctx) error {
	var request dto.CreateStoreSalePaymentRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[CreateStoreSalePayment] failed to parse request", zap.Error(err))
		return err
	}

	storeSaleIdParam := c.Params("storeSaleId")
	if storeSaleIdParam == "" {
		h.log.Error("[CreateStoreSalePayment] storeSaleId is required")
		return errx.BadRequest("storeSaleId is required")
	}

	storeSaleId, err := strconv.ParseUint(storeSaleIdParam, 10, 64)
	if err != nil {
		h.log.Error("[CreateStoreSalePayment] failed to parse storeSaleId", zap.Error(err))
		return errx.BadRequest("failed to parse storeSaleId")
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[CreateStoreSalePayment] failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[CreateStoreSalePayment] failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	res, err := h.service.CreateStoreSalePayment(storeSaleId, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create store sale payment")
}

func (h *StoreHandler) UpdateStoreSale(c *fiber.Ctx) error {
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

	var request dto.UpdateStoreSaleRequest
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

	res, err := h.service.UpdateStoreSale(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update store sale")
}

func (h *StoreHandler) UpdateStoreSalePayment(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[UpdateStoreSalePayment] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("[UpdateStoreSalePayment] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	var request dto.UpdateStoreSalePaymentRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[UpdateStoreSalePayment] failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[UpdateStoreSalePayment] failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[UpdateStoreSalePayment] failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	res, err := h.service.UpdateStoreSalePayment(id, request, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[UpdateStoreSalePayment] failed to update store sale payment", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update store sale payment")
}

func (h *StoreHandler) SendStoreSale(c *fiber.Ctx) error {
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

	res, err := h.service.SendStoreSale(storeSaleId, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[SendStoreSale] failed to send store sale", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success send store sale")
}
