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

	v1.Post("/sales", middleware.Authentication(), h.CreateStoreSale)
	v1.Get("/sales/:id", middleware.Authentication(), h.GetStoreSaleById)
	v1.Get("/sales", middleware.Authentication(), h.GetStoreSales)
	v1.Put("/sales/:id", middleware.Authentication(), h.UpdateStoreSale)
	v1.Delete("/sales/:id", middleware.Authentication(), h.DeleteStoreSale)
	v1.Post("/sales/:storeSaleId/payments", middleware.Authentication(), h.CreateStoreSalePayment)
	v1.Put("/sales/:storeSaleId/payments/:id", middleware.Authentication(), h.UpdateStoreSalePayment)
	v1.Delete("/sales/:storeSaleId/payments/:id", middleware.Authentication(), h.DeleteStoreSalePayment)
	v1.Patch("sales/:storeSaleId/send", middleware.Authentication(), h.SendStoreSale)

	v1.Get("/", middleware.Authentication(), h.GetStores)
	v1.Post("/", middleware.Authentication(), h.CreateStore)
	v1.Put("/:id", middleware.Authentication(), h.UpdateStore)
	v1.Delete("/:id", middleware.Authentication(), h.DeleteStore)
	v1.Get("/:id", middleware.Authentication(), h.GetStoreDetail)
	v1.Get("/overview/stocks/:id", middleware.Authentication(), h.GetStoreItemStocks)
	v1.Get("/overview", middleware.Authentication(), h.GetStoreOverview)

	v1.Post("/request/items", middleware.Authentication(), h.CreateStoreRequestItem)
	v1.Get("/request/items", middleware.Authentication(), h.GetStockRequestItems)
	v1.Get("/request/items/:id", middleware.Authentication(), h.GetStoreRequestItemById)
	v1.Put("/request/items/:id", middleware.Authentication(), h.UpdateStoreRequestItem)
	v1.Put("/request/items/:id/warehouse-confirmations", middleware.Authentication(), h.WarehouseConfirmationStoreRequestItem)
	v1.Put("/request/items/:id/store-confirmations", middleware.Authentication(), h.StoreConfirmationStoreRequestItem)
	v1.Put("/request/items/:id/sorting-cracked-eggs", middleware.Authentication(), h.SortingStoreRequestItem)

	v1.Get("/:storeId/items/:itemId", middleware.Authentication(), h.GetStoreItem)
	v1.Put("/:storeId/items/:itemId", middleware.Authentication(), h.UpdateStoreItem)
	v1.Get("/items/eggs/summary/:storeId", middleware.Authentication(), h.GetEggStoreItemSummary)

	v1.Get("/items/histories", middleware.Authentication(), h.GetStoreItemHistories)
	v1.Get("/items/histories/:id", middleware.Authentication(), h.GetStoreItemHistory)
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
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	res, err := h.service.GetStoreRequestItems(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get stock request items")
}

func (h *StoreHandler) GetStoreRequestItemById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	res, err := h.service.GetStoreRequestItemById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get store request item by id")
}

func (h *StoreHandler) WarehouseConfirmationStoreRequestItem(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	var request dto.WarehouseConfirmationStoreRequestItem
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

	res, err := h.service.WarehouseConfirmationStoreRequestItem(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success confirm store request item by warehouse")
}

func (h *StoreHandler) StoreConfirmationStoreRequestItem(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	var request dto.StoreConfirmationStoreRequestItem
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

	res, err := h.service.StoreConfirmationStoreRequestItem(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success confirm store request item by store")
}

func (h *StoreHandler) UpdateStoreRequestItem(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	var request dto.UpdateStoreRequestItemRequest
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
		h.log.Error("failed to get user id from context")
		return errx.Unauthorized("no user id in context")
	}

	res, err := h.service.UpdateStoreRequestItem(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update request store item")
}

func (h *StoreHandler) SortingStoreRequestItem(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	var request dto.SortingStoreRequestItemRequest
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
		h.log.Error("failed to get user id from context")
		return errx.Unauthorized("no user id in context")
	}

	res, err := h.service.SortingStoreRequestItem(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update request store item")
}

func (h *StoreHandler) GetStoreItemStocks(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	data, err := h.service.GetStoreItemStocks(id)
	if err != nil {
		h.log.Error("failed to get store item overview", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get store item overview")
}

func (h *StoreHandler) GetStoreOverview(c *fiber.Ctx) error {
	var filter dto.GetStoreOverviewFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&filter); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return err
	}

	data, err := h.service.GetStoreOverview(filter)
	if err != nil {
		h.log.Error("failed to get store overview", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get store overview")
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

func (h *StoreHandler) GetEggStoreItemSummary(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("storeId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse store id")
		return errx.BadRequest("invalid store id")
	}

	data, err := h.service.GetEggStoreItemSummary(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get egg warehouse item summary")
}

func (h *StoreHandler) GetStoreItemHistories(c *fiber.Ctx) error {
	var filter dto.GetStoreItemHistoryFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query param", zap.Error(err))
		return err
	}

	data, err := h.service.GetStoreItemHistories(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get store item histories")
}

func (h *StoreHandler) GetStoreItemHistory(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	data, err := h.service.GetStoreItemHistoryById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get store item history")
}

func (h *StoreHandler) CreateStoreSale(c *fiber.Ctx) error {
	var request dto.CreateStoreSaleRequest
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

	res, err := h.service.CreateStoreSale(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create store sale")
}

func (h *StoreHandler) GetStoreSaleById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	res, err := h.service.GetStoreSaleById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get store sale by id")
}

func (h *StoreHandler) DeleteStoreSale(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	err = h.service.DeleteStoreSale(id, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *StoreHandler) GetStoreSales(c *fiber.Ctx) error {
	var filter dto.GetStoreSaleFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
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
		h.log.Error("failed to parse request", zap.Error(err))
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

	res, err := h.service.CreateStoreSalePayment(storeSaleId, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create store sale payment")
}

func (h *StoreHandler) UpdateStoreSale(c *fiber.Ctx) error {
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

	var request dto.UpdateStoreSaleRequest
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

	res, err := h.service.UpdateStoreSale(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update store sale")
}

func (h *StoreHandler) UpdateStoreSalePayment(c *fiber.Ctx) error {
	storeSaleId, err := strconv.ParseUint(c.Params("storeSaleId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse store sale id", zap.Error(err))
		return errx.BadRequest("failed to store sale id")
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	var request dto.UpdateStoreSalePaymentRequest
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

	res, err := h.service.UpdateStoreSalePayment(storeSaleId, id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update store sale payment")
}

func (h *StoreHandler) DeleteStoreSalePayment(c *fiber.Ctx) error {
	storeSaleId, err := strconv.ParseUint(c.Params("storeSaleId"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse store sale id", zap.Error(err))
		return errx.BadRequest("failed to store sale id")
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

	err = h.service.DeleteStoreSalePayment(storeSaleId, id, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *StoreHandler) SendStoreSale(c *fiber.Ctx) error {
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

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	res, err := h.service.SendStoreSale(storeSaleId, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success send store sale")
}
