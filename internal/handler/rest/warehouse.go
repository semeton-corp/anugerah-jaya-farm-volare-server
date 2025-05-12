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
	v1.Get("/", middleware.Authentication(), a.GetWarehouses)

	v1.Post("/items", middleware.Authentication(), a.CreateWarehouseItem)
	v1.Get("/items", middleware.Authentication(), a.GetWarehouseItem)
	v1.Get("/items/:id", middleware.Authentication(), a.GetWarehouseItemById)
	v1.Put("/items/:id", middleware.Authentication(), a.UpdateWarehouseItem)
	v1.Delete("/items/:id", middleware.Authentication(), a.DeleteWarehouseItem)

	v1.Post("/stock-items", middleware.Authentication(), a.CreateWarehouseStockItem)
	v1.Get("/stock-items", middleware.Authentication(), a.GetWarehouseStockItems)
	v1.Get("/:warehouseId/stock-items/:warehouseItemId", middleware.Authentication(), a.GetWarehouseStockItemByWarehouseIdAndWarehouseItemId)
	v1.Put("/:warehouseId/stock-items/:warehouseItemId", middleware.Authentication(), a.UpdateWarehouseStockItem)
	v1.Delete("/:warehouseId/stock-items/:warehouseItemId", middleware.Authentication(), a.DeleteWarehouseStockItem)

	v1.Post("/order-items", middleware.Authentication(), a.CreateWarehouseOrderItem)
	v1.Get("/order-items", middleware.Authentication(), a.GetWarehouseOrderItems)
	v1.Get("/order-items/:id", middleware.Authentication(), a.GetWarehouseOrderItemById)
	v1.Delete("/order-items/:id", middleware.Authentication(), a.DeleteWarehouseOrderItem)
	v1.Patch("/order-items/:id/takes", middleware.Authentication(), a.TakeWarehouseOrderItem)

	v1.Post("/items/convert/good-egg/butir-to-ikat", middleware.Authentication(), a.GoodEggConvertionButirToIkat)
	v1.Post("/items/convert/good-egg/ikat-to-butir", middleware.Authentication(), a.GoodEggConvertionIkatToButir)
	v1.Post("/items/convert/cracked-egg/butir-to-pack", middleware.Authentication(), a.CrackedEggConverterButirToPack)
}

func NewWarehouseHandler(log *zap.Logger, service service.IWarehouseService, validator *validator.Validate) *WarehouseHandler {
	return &WarehouseHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (a *WarehouseHandler) GetWarehouses(c *fiber.Ctx) error {
	warehouses, err := a.service.GetWarehouses()
	if err != nil {
		a.log.Error("[GetWarehouses] failed to get warehouses", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, warehouses, "get warehouses success")
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

	res, err := a.service.CreateWarehouseItem(request, uuid.MustParse(accountId))
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

	warehouseItems, err := a.service.GetWarehouseItems(filter)
	if err != nil {
		a.log.Error("[GetWarehouseItem] failed to get warehouse items", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, warehouseItems, "get warehouse items success")
}

func (a *WarehouseHandler) GetWarehouseItemById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		a.log.Error("[GetWarehouseItemById] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		a.log.Error("[GetWarehouseItemById] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	res, err := a.service.GetWarehouseItemById(id)
	if err != nil {
		a.log.Error("[GetWarehouseItemById] failed to get warehouse item", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "get warehouse item success")
}

func (a *WarehouseHandler) UpdateWarehouseItem(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		a.log.Error("[UpdateWarehouseItem] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		a.log.Error("[UpdateWarehouseItem] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	var request dto.UpdateWarehouseItemRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[UpdateWarehouseItem] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[UpdateWarehouseItem] failed to validate request", zap.Error(err))
		return err
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[UpdateWarehouseItem] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	res, err := a.service.UpdateWarehouseItem(id, request, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[UpdateWarehouseItem] failed to update warehouse item", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "update warehouse item success")
}

func (a *WarehouseHandler) DeleteWarehouseItem(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		a.log.Error("[DeleteWarehouseItem] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		a.log.Error("[DeleteWarehouseItem] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	err = a.service.DeleteWarehouseItem(id)
	if err != nil {
		a.log.Error("[DeleteWarehouseItem] failed to delete warehouse item", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}

func (a *WarehouseHandler) CreateWarehouseStockItem(c *fiber.Ctx) error {
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

func (a *WarehouseHandler) GetWarehouseStockItems(c *fiber.Ctx) error {
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

func (a *WarehouseHandler) GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(c *fiber.Ctx) error {
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

func (a *WarehouseHandler) UpdateWarehouseStockItem(c *fiber.Ctx) error {
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

func (a *WarehouseHandler) DeleteWarehouseStockItem(c *fiber.Ctx) error {
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

func (a *WarehouseHandler) CreateWarehouseOrderItem(c *fiber.Ctx) error {
	var request dto.CreateWarehouseOrderItemRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[CreateStoreOrderItem] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[CreateStoreOrderItem] failed to validate request", zap.Error(err))
		return err
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[CreateStoreOrderItem] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	res, err := a.service.CreateWarehouseOrderItem(request, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[CreateStoreOrderItem] failed to create store order item", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "create store order item success")
}

func (a *WarehouseHandler) GetWarehouseOrderItemById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		a.log.Error("[GetWarehouseOrderItemById] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		a.log.Error("[GetWarehouseOrderItemById] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	res, err := a.service.GetWarehouseOrderItemById(id)
	if err != nil {
		a.log.Error("[GetWarehouseOrderItemById] failed to get store order item", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "get store order item success")
}

func (a *WarehouseHandler) GetWarehouseOrderItems(c *fiber.Ctx) error {
	var filter dto.GetWarehouseOrderItemFilter
	if err := c.QueryParser(&filter); err != nil {
		a.log.Error("[GetWarehouseOrderItems] failed to parse query", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(filter); err != nil {
		a.log.Error("[GetWarehouseOrderItems] failed to validate request", zap.Error(err))
		return err
	}

	warehouseOrderItems, err := a.service.GetWarehouseOrderItems(filter)
	if err != nil {
		a.log.Error("[GetWarehouseOrderItems] failed to get store order items", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, warehouseOrderItems, "get store order items success")
}

func (a *WarehouseHandler) DeleteWarehouseOrderItem(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		a.log.Error("[DeleteWarehouseOrderItem] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		a.log.Error("[DeleteWarehouseOrderItem] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	err = a.service.DeleteWarehouseOrderItem(id)
	if err != nil {
		a.log.Error("[DeleteWarehouseOrderItem] failed to delete store order item", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}

func (a *WarehouseHandler) TakeWarehouseOrderItem(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		a.log.Error("[TakeWarehouseOrderItem] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		a.log.Error("[TakeWarehouseOrderItem] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[TakeWarehouseOrderItem] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	res, err := a.service.TakeWarehouseOrderItem(id, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[TakeWarehouseOrderItem] failed to take store order item", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "take store order item success")
}

func (a *WarehouseHandler) GoodEggConvertionButirToIkat(c *fiber.Ctx) error {
	var request dto.GoodEggWarehouseConvertionRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[GoodEggConvertionButirToIkat] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[GoodEggConvertionButirToIkat] failed to validate request", zap.Error(err))
		return err
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[GoodEggConvertionButirToIkat] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	resp, err := a.service.GoodEggConvertionButirToIkat(request, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[GoodEggConvertionButirToIkat] failed to good egg convertion butir to ikat", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "good egg convertion butir to ikat success")
}

func (a *WarehouseHandler) GoodEggConvertionIkatToButir(c *fiber.Ctx) error {
	var request dto.GoodEggWarehouseConvertionRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[GoodEggConvertionIkatToButir] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[GoodEggConvertionIkatToButir] failed to validate request", zap.Error(err))
		return err
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[GoodEggConvertionIkatToButir] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	resp, err := a.service.GoodEggConvertionIkatToButir(request, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[GoodEggConvertionIkatToButir] failed to good egg convertion ikat to butir", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "good egg convertion ikat to butir success")
}

func (a *WarehouseHandler) CrackedEggConverterButirToPack(c *fiber.Ctx) error {
	var request dto.CrackedEggWarehouseConvertionRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[CrackedEggConverterButirToPack] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[CrackedEggConverterButirToPack] failed to validate request", zap.Error(err))
		return err
	}

	accountId, ok := c.Locals("accountId").(string)
	if !ok {
		a.log.Error("[CrackedEggConverterButirToPack] failed to get accountId from context")
		return errx.Unauthorized("no accountId in context")
	}

	resp, err := a.service.CrackedEggConvertionButirToPack(request, uuid.MustParse(accountId))
	if err != nil {
		a.log.Error("[CrackedEggConverterButirToPacket] failed to cracked egg converter butir to packet", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "cracked egg converter butir to packet success")
}
