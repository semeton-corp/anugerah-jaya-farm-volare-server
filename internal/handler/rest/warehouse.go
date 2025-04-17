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

type WarehouseHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IWarehouseService
}

func (a *WarehouseHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/warehouses")
	v1.Post("/items", middleware.Authentication(), a.CreateWarehouseItem)
	v1.Get("/items", middleware.Authentication(), a.GetWarehouseItem)

	v1.Post("/stock-items", middleware.Authentication(), a.CreateStockWarehouseItem)
	v1.Get("/stock-items", middleware.Authentication(), a.GetStockWarehouseItems)
	v1.Get("/:warehouseId/stock-items/:warehouseItemId", middleware.Authentication(), a.GetStockWarehouseItemByWarehouseIdAndWarehouseItemId)
	v1.Put("/:warehouseId/stock-items/:warehouseItemId", middleware.Authentication(), a.UpdateStockWarehouseItem)
	v1.Delete("/:warehouseId/stock-items/:warehouseItemId", middleware.Authentication(), a.DeleteStockWarehouseItem)
}

func NewWarehouseHandler(log *zap.Logger, service service.IWarehouseService, validator *validator.Validate) *WarehouseHandler {
	return &WarehouseHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (a *WarehouseHandler) CreateWarehouseItem(c *fiber.Ctx) error {
	var request dto.CreateWarehouseItemRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[CreateWarehouseItem] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[CreateWarehouseItem] failed to validate request", zap.Error(err))
		return err
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[CreateWarehouseItem] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	res, err := a.service.CreateWarehouseItem(&request, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[CreateWarehouseItem] failed to create warehouse item", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "create warehouse item success")
}

func (a *WarehouseHandler) GetWarehouseItem(c *fiber.Ctx) error {
	var filter dto.GetWarehouseItemFilter
	if err := c.QueryParser(&filter); err != nil {
		a.log.Error("[GetWarehouseItem] failed to parse query", zap.Error(err))
		return err
	}

	warehouseItems, err := a.service.GetWarehouseItem(filter)
	if err != nil {
		a.log.Error("[GetWarehouseItem] failed to get warehouse items", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, warehouseItems, "get warehouse items success")
}

func (a *WarehouseHandler) CreateStockWarehouseItem(c *fiber.Ctx) error {
	var request dto.CreateWarehouseStockItemRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[CreateStockWarehouseItem] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[CreateStockWarehouseItem] failed to validate request", zap.Error(err))
		return err
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[CreateStockWarehouseItem] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	res, err := a.service.CreateWarehouseStockItem(&request, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[CreateStockWarehouseItem] failed to create stock warehouse item", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "create stock warehouse item success")
}

func (a *WarehouseHandler) GetStockWarehouseItems(c *fiber.Ctx) error {
	var filter dto.GetWarehouseStockItemFilter
	if err := c.QueryParser(&filter); err != nil {
		a.log.Error("[GetStockWarehouseItems] failed to parse query", zap.Error(err))
		return err
	}

	stockWarehouseItems, err := a.service.GetWarehouseStockItems(filter)
	if err != nil {
		a.log.Error("[GetStockWarehouseItems] failed to get stock warehouse items", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, stockWarehouseItems, "get stock warehouse items success")
}

func (a *WarehouseHandler) GetStockWarehouseItemByWarehouseIdAndWarehouseItemId(c *fiber.Ctx) error {
	warehouseIdParam := c.Params("warehouseId")
	warehouseItemIdParam := c.Params("warehouseItemId")

	if warehouseIdParam == "" || warehouseItemIdParam == "" {
		a.log.Error("[GetStockWarehouseItemByWarehouseIdAndWarehouseItemId] warehouseId and warehouseItemId are required")
		return errx.BadRequest("warehouseId and warehouseItemId are required")
	}

	warehouseId, err := strconv.ParseUint(warehouseIdParam, 10, 64)
	if err != nil {
		a.log.Error("[GetStockWarehouseItemByWarehouseIdAndWarehouseItemId] failed to parse warehouseId", zap.Error(err))
		return errx.BadRequest("failed to parse warehouseId")
	}

	warehouseItemId, err := strconv.ParseUint(warehouseItemIdParam, 10, 64)
	if err != nil {
		a.log.Error("[GetStockWarehouseItemByWarehouseIdAndWarehouseItemId] failed to parse warehouseItemId", zap.Error(err))
		return errx.BadRequest("failed to parse warehouseItemId")
	}

	res, err := a.service.GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId, warehouseItemId)
	if err != nil {
		a.log.Error("[GetStockWarehouseItemByWarehouseIdAndWarehouseItemId] failed to get stock warehouse item", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "get stock warehouse item success")
}

func (a *WarehouseHandler) UpdateStockWarehouseItem(c *fiber.Ctx) error {
	warehouseIdParam := c.Params("warehouseId")
	warehouseItemIdParam := c.Params("warehouseItemId")

	if warehouseIdParam == "" || warehouseItemIdParam == "" {
		a.log.Error("[GetStockWarehouseItemByWarehouseIdAndWarehouseItemId] warehouseId and warehouseItemId are required")
		return errx.BadRequest("warehouseId and warehouseItemId are required")
	}

	warehouseId, err := strconv.ParseUint(warehouseIdParam, 10, 64)
	if err != nil {
		a.log.Error("[GetStockWarehouseItemByWarehouseIdAndWarehouseItemId] failed to parse warehouseId", zap.Error(err))
		return errx.BadRequest("failed to parse warehouseId")
	}

	warehouseItemId, err := strconv.ParseUint(warehouseItemIdParam, 10, 64)
	if err != nil {
		a.log.Error("[GetStockWarehouseItemByWarehouseIdAndWarehouseItemId] failed to parse warehouseItemId", zap.Error(err))
		return errx.BadRequest("failed to parse warehouseItemId")
	}

	var request dto.UpdateWarehouseStockItemRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[UpdateStockWarehouseItem] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[UpdateStockWarehouseItem] failed to validate request", zap.Error(err))
		return err
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[UpdateStockWarehouseItem] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	res, err := a.service.UpdateWarehouseStockItem(warehouseId, warehouseItemId, request, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[UpdateStockWarehouseItem] failed to update stock warehouse item", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "update stock warehouse item success")
}

func (a *WarehouseHandler) DeleteStockWarehouseItem(c *fiber.Ctx) error {
	warehouseIdParam := c.Params("warehouseId")
	warehouseItemIdParam := c.Params("warehouseItemId")

	if warehouseIdParam == "" || warehouseItemIdParam == "" {
		a.log.Error("[DeleteStockWarehouseItem] warehouseId and warehouseItemId are required")
		return errx.BadRequest("warehouseId and warehouseItemId are required")
	}

	warehouseId, err := strconv.ParseUint(warehouseIdParam, 10, 64)
	if err != nil {
		a.log.Error("[DeleteStockWarehouseItem] failed to parse warehouseId", zap.Error(err))
		return errx.BadRequest("failed to parse warehouseId")
	}

	warehouseItemId, err := strconv.ParseUint(warehouseItemIdParam, 10, 64)
	if err != nil {
		a.log.Error("[DeleteStockWarehouseItem] failed to parse warehouseItemId", zap.Error(err))
		return errx.BadRequest("failed to parse warehouseItemId")
	}

	err = a.service.DeleteWarehouseStockItem(warehouseId, warehouseItemId)
	if err != nil {
		a.log.Error("[DeleteStockWarehouseItem] failed to delete stock warehouse item", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}
