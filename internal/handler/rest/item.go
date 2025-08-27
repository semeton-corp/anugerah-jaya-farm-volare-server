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

type EggPriceHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.IItemService
}

func (h *EggPriceHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/items")
	v1.Post("/prices/discounts", middleware.Authentication(), h.CreateItemPriceDiscount)
	v1.Get("/prices/discounts", middleware.Authentication(), h.GetItemPriceDiscounts)
	v1.Get("/prices/discounts/:id", middleware.Authentication(), h.GetItemPriceDiscountById)
	v1.Put("/prices/discounts/:id", middleware.Authentication(), h.UpdateItemPriceDiscount)
	v1.Delete("/prices/discounts/:id", middleware.Authentication(), h.DeleteItemPriceDiscount)

	v1.Post("/prices", middleware.Authentication(), h.CreateItemPrice)
	v1.Get("/prices", middleware.Authentication(), h.GetItemPrices)
	v1.Get("/prices/:id", middleware.Authentication(), h.GetItemPriceById)
	v1.Put("/prices/:id", middleware.Authentication(), h.UpdateItemPrice)
	v1.Delete("/prices/:id", middleware.Authentication(), h.DeleteItemPrice)

	v1.Post("/", middleware.Authentication(), h.CreateItem)
	v1.Get("/", middleware.Authentication(), h.GetItems)
	v1.Get("/:id", middleware.Authentication(), h.GetItemById)
	v1.Put("/:id", middleware.Authentication(), h.UpdateItem)
	v1.Delete("/:id", middleware.Authentication(), h.DeleteItem)
}

func NewEggPriceHandler(log *zap.Logger, service service.IItemService, validator *validator.Validate) *EggPriceHandler {
	return &EggPriceHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *EggPriceHandler) CreateItemPrice(c *fiber.Ctx) error {
	var request dto.CreateItemPriceRequest
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
		h.log.Error("userId not found in context")
		return errx.Unauthorized("userId not found in context")
	}

	resp, err := h.service.CreateItemPrice(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusCreated,
		resp,
		"success create item price",
	)
}

func (h *EggPriceHandler) GetItemPrices(c *fiber.Ctx) error {
	eggPrices, err := h.service.GetItemPrices()
	if err != nil {
		h.log.Error("failed to get item prices", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		eggPrices,
		"success get item prices",
	)
}

func (h *EggPriceHandler) GetItemPriceById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	resp, err := h.service.GetItemPriceById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		resp,
		"success get item price by id",
	)
}

func (h *EggPriceHandler) UpdateItemPrice(c *fiber.Ctx) error {
	var request dto.UpdateItemPriceRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("userId not found in context")
		return errx.Unauthorized("userId not found in context")
	}

	resp, err := h.service.UpdateItemPrice(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		resp,
		"success update item price",
	)
}

func (h *EggPriceHandler) DeleteItemPrice(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return err
	}

	err = h.service.DeleteItemPrice(id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *EggPriceHandler) CreateItemPriceDiscount(c *fiber.Ctx) error {
	var filter dto.CreateItemPriceDiscountRequest
	if err := c.BodyParser(&filter); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("userId not found in context")
		return errx.Unauthorized("userId not found in context")
	}

	resp, err := h.service.CreateItemDiscount(filter, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusCreated,
		resp,
		"success create item price discount",
	)
}

func (h *EggPriceHandler) GetItemPriceDiscountById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return err
	}

	resp, err := h.service.GetItemDiscountById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		resp,
		"success get item price discount by id",
	)
}

func (h *EggPriceHandler) GetItemPriceDiscounts(c *fiber.Ctx) error {
	resp, err := h.service.GetItemDiscounts()
	if err != nil {
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		resp,
		"success get item price discounts",
	)
}

func (h *EggPriceHandler) UpdateItemPriceDiscount(c *fiber.Ctx) error {
	var request dto.UpdateItemPriceDiscountRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("userId not found in context")
		return errx.Unauthorized("userId not found in context")
	}

	resp, err := h.service.UpdateItemDiscount(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		resp,
		"success update item price discount",
	)
}

func (h *EggPriceHandler) DeleteItemPriceDiscount(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return err
	}

	err = h.service.DeleteItemDiscount(id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}

func (h *EggPriceHandler) CreateItem(c *fiber.Ctx) error {
	var request dto.CreateItemRequest
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

	res, err := h.service.CreateItem(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, res, "create item success")
}

func (h *EggPriceHandler) GetItems(c *fiber.Ctx) error {
	var filter dto.GetItemFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query", zap.Error(err))
		return err
	}

	Items, err := h.service.GetItems(filter)
	if err != nil {
		h.log.Error("failed to get items", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, Items, "get items success")
}

func (h *EggPriceHandler) GetItemById(c *fiber.Ctx) error {
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

	res, err := h.service.GetItemById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "get item success")
}

func (h *EggPriceHandler) UpdateItem(c *fiber.Ctx) error {
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

	var request dto.UpdateItemRequest
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

	res, err := h.service.UpdateItem(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, res, "update item success")
}

func (h *EggPriceHandler) DeleteItem(c *fiber.Ctx) error {
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

	err = h.service.DeleteItem(id)
	if err != nil {
		return err
	}

	return response.NoContentResponse(c)
}
