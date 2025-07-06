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
	v1.Post("/", middleware.Authentication(), h.CreateWarehouse)
	v1.Put("/:id", middleware.Authentication(), h.UpdateWarehouse)
	v1.Delete("/:id", middleware.Authentication(), h.DeleteWarehouse)
	v1.Get("/", middleware.Authentication(), h.GetWarehouses)
	v1.Get("/:id", middleware.Authentication(), h.GetWarehouseDetail)

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

	v1.Post("/items/convert/good-egg/butir-to-ikat", middleware.Authentication(), h.GoodEggConvertionButirToIkat)
	v1.Post("/items/convert/good-egg/ikat-to-butir", middleware.Authentication(), h.GoodEggConvertionIkatToButir)
	v1.Post("/items/convert/cracked-egg/butir-to-pack", middleware.Authentication(), h.CrackedEggConverterButirToPack)
}

func NewWarehouseHandler(log *zap.Logger, service service.IWarehouseService, validator *validator.Validate) *WarehouseHandler {
	return &WarehouseHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
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

	return response.SuccessResponse(c, fiber.StatusCreated, res, "create stock warehouse item success")
}

func (h *WarehouseHandler) GetWarehouseItems(c *fiber.Ctx) error {
	var filter dto.GetWarehouseItemFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	stockWarehouseItems, err := h.service.GetWarehouseItems(filter)
	if err != nil {
		h.log.Error("failed to get stock warehouse items", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, stockWarehouseItems, "get stock warehouse items success")
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

	return response.SuccessResponse(c, fiber.StatusOK, res, "get stock warehouse item success")
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
		h.log.Error("failed to update stock warehouse item", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "update stock warehouse item success")
}

func (h *WarehouseHandler) DeleteWarehouseItem(c *fiber.Ctx) error {
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

	err = h.service.DeleteWarehouseItem(warehouseId, warehouseItemId)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *WarehouseHandler) CreateWarehouseOrderItem(c *fiber.Ctx) error {
	var request dto.CreateWarehouseOrderItemRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[CreateStoreOrderItem] failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[CreateStoreOrderItem] failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[CreateStoreOrderItem] failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	res, err := h.service.CreateWarehouseOrderItem(request, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[CreateStoreOrderItem] failed to create store order item", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "create store order item success")
}

func (h *WarehouseHandler) GetWarehouseOrderItemById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[GetWarehouseOrderItemById] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("[GetWarehouseOrderItemById] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	res, err := h.service.GetWarehouseOrderItemById(id)
	if err != nil {
		h.log.Error("[GetWarehouseOrderItemById] failed to get store order item", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "get store order item success")
}

func (h *WarehouseHandler) GetWarehouseOrderItems(c *fiber.Ctx) error {
	var filter dto.GetWarehouseOrderItemFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("[GetWarehouseOrderItems] failed to parse query", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("[GetWarehouseOrderItems] failed to validate request", zap.Error(err))
		return err
	}

	warehouseOrderItems, err := h.service.GetWarehouseOrderItems(filter)
	if err != nil {
		h.log.Error("[GetWarehouseOrderItems] failed to get store order items", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, warehouseOrderItems, "get store order items success")
}

func (h *WarehouseHandler) DeleteWarehouseOrderItem(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[DeleteWarehouseOrderItem] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("[DeleteWarehouseOrderItem] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	err = h.service.DeleteWarehouseOrderItem(id)
	if err != nil {
		h.log.Error("[DeleteWarehouseOrderItem] failed to delete store order item", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}

func (h *WarehouseHandler) TakeWarehouseOrderItem(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		h.log.Error("[TakeWarehouseOrderItem] id is required")
		return errx.BadRequest("id is required")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Error("[TakeWarehouseOrderItem] failed to parse id", zap.Error(err))
		return errx.BadRequest("failed to parse id")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[TakeWarehouseOrderItem] failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	res, err := h.service.TakeWarehouseOrderItem(id, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[TakeWarehouseOrderItem] failed to take store order item", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "take store order item success")
}

func (h *WarehouseHandler) GoodEggConvertionButirToIkat(c *fiber.Ctx) error {
	var request dto.GoodEggWarehouseConvertionRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[GoodEggConvertionButirToIkat] failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[GoodEggConvertionButirToIkat] failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[GoodEggConvertionButirToIkat] failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	resp, err := h.service.GoodEggConvertionButirToIkat(request, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[GoodEggConvertionButirToIkat] failed to good egg convertion butir to ikat", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "good egg convertion butir to ikat success")
}

func (h *WarehouseHandler) GoodEggConvertionIkatToButir(c *fiber.Ctx) error {
	var request dto.GoodEggWarehouseConvertionRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[GoodEggConvertionIkatToButir] failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[GoodEggConvertionIkatToButir] failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[GoodEggConvertionIkatToButir] failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	resp, err := h.service.GoodEggConvertionIkatToButir(request, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[GoodEggConvertionIkatToButir] failed to good egg convertion ikat to butir", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "good egg convertion ikat to butir success")
}

func (h *WarehouseHandler) CrackedEggConverterButirToPack(c *fiber.Ctx) error {
	var request dto.CrackedEggWarehouseConvertionRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[CrackedEggConverterButirToPack] failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[CrackedEggConverterButirToPack] failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[CrackedEggConverterButirToPack] failed to get userId from context")
		return errx.Unauthorized("no userId in context")
	}

	resp, err := h.service.CrackedEggConvertionButirToPack(request, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[CrackedEggConverterButirToPacket] failed to cracked egg converter butir to packet", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "cracked egg converter butir to packet success")
}
