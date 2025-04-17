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

func (a *StoreHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/stores")
	v1.Get("/", middleware.Authentication(), a.GetStores)

	v1.Post("/request-items", middleware.Authentication(), a.CreateStoreRequestItem)
	v1.Get("/request-items", middleware.Authentication(), a.GetStockRequestItems)
	v1.Get("/request-items/:id", middleware.Authentication(), a.GetStoreRequestItemById)
	v1.Put("/request-items/:id/warehouses", middleware.Authentication(), a.UpdateStoreRequestItemByWarehouse)
	v1.Put("/request-items/:id/stores", middleware.Authentication(), a.UpdateStoreRequestItemByStore)
}

func NewStoreHandler(log *zap.Logger, service service.IStoreService, validator *validator.Validate) *StoreHandler {
	return &StoreHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (a *StoreHandler) GetStores(c *fiber.Ctx) error {
	stores, err := a.service.GetStores()
	if err != nil {
		a.log.Error("[GetStores] failed to get stores", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, stores, "success get stores")
}

func (a *StoreHandler) CreateStoreRequestItem(c *fiber.Ctx) error {
	var request dto.CreateStoreRequestItemRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[CreateStoreRequestItem] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[CreateStoreRequestItem] failed to validate request", zap.Error(err))
		return err
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[CreateStoreRequestItem] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	res, err := a.service.CreateStoreRequestItem(request, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[CreateStoreRequestItem] failed to create store request item", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "success create store request item")
}

func (a *StoreHandler) GetStockRequestItems(c *fiber.Ctx) error {
	var filter dto.GetStoreRequestItemFilter
	if err := c.QueryParser(&filter); err != nil {
		a.log.Error("[GetStockRequestItems] failed to parse query", zap.Error(err))
		return err
	}

	res, err := a.service.GetStoreRequestItems(filter)
	if err != nil {
		a.log.Error("[GetStockRequestItems] failed to get stock request items", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get stock request items")
}

func (a *StoreHandler) GetStoreRequestItemById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		a.log.Error("[GetStoreRequestItemById] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		a.log.Error("[GetStoreRequestItemById] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	res, err := a.service.GetStoreRequestItemById(id)
	if err != nil {
		a.log.Error("[GetStoreRequestItemById] failed to get store request item by id", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success get store request item by id")
}

func (a *StoreHandler) UpdateStoreRequestItemByWarehouse(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		a.log.Error("[UpdateStoreRequestItemByWarehouse] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		a.log.Error("[UpdateStoreRequestItemByWarehouse] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	var request dto.UpdateStoreRequestItemByWarehouseRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[UpdateStoreRequestItemByWarehouse] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[UpdateStoreRequestItemByWarehouse] failed to validate request", zap.Error(err))
		return err
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[UpdateStoreRequestItemByWarehouse] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	storeRequestItem, err := a.service.GetStoreRequestItemById(id)
	if err != nil {
		a.log.Error("[UpdateStoreRequestItemByWarehouse] failed to get store request item by id", zap.Error(err))
		return err
	}

	combinedRequest := dto.UpdateStoreRequestItemRequest{
		Quantity: storeRequestItem.Quantity,
		Status:   request.Status,
	}

	res, err := a.service.UpdateStoreRequestItem(id, combinedRequest, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[UpdateStoreRequestItemByWarehouse] failed to update store request item by warehouse", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update store request item by warehouse")
}

func (a *StoreHandler) UpdateStoreRequestItemByStore(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		a.log.Error("[UpdateStoreRequestItemByWarehouse] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		a.log.Error("[UpdateStoreRequestItemByWarehouse] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	var request dto.UpdateStoreRequestItemByStoreRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[UpdateStoreRequestItemByWarehouse] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[UpdateStoreRequestItemByWarehouse] failed to validate request", zap.Error(err))
		return err
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[UpdateStoreRequestItemByWarehouse] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	// storeRequestItem, err := a.service.GetStoreRequestItemById(id)
	// if err != nil {
	// 	a.log.Error("[UpdateStoreRequestItemByWarehouse] failed to get store request item by id", zap.Error(err))
	// 	return err
	// }

	combinedRequest := dto.UpdateStoreRequestItemRequest{
		Quantity: request.Quantity,
		Status:   request.Status,
	}

	res, err := a.service.UpdateStoreRequestItem(id, combinedRequest, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[UpdateStoreRequestItemByWarehouse] failed to update store request item by warehouse", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "success update store request item by warehouse")
}
