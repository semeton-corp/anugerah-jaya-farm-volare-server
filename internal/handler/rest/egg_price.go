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
	service   service.IEggPriceService
}

func (h *EggPriceHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/egg-prices")
	v1.Post("", middleware.Authentication(), h.CreateEggPrice)
	v1.Get("", middleware.Authentication(), h.GetEggPrices)
	v1.Get("/:id", middleware.Authentication(), h.GetEggPriceById)
	v1.Put("/:id", middleware.Authentication(), h.UpdateEggPrice)
	v1.Delete("/:id", middleware.Authentication(), h.DeleteEggPrice)

	v1.Post("/discounts", middleware.Authentication(), h.CreateEggPriceDiscount)
	v1.Get("/discounts", middleware.Authentication(), h.GetEggPriceDiscounts)
	v1.Get("/discounts/:id", middleware.Authentication(), h.GetEggPriceDiscountById)
	v1.Put("/discounts/:id", middleware.Authentication(), h.UpdateEggPriceDiscount)
	v1.Delete("/discounts/:id", middleware.Authentication(), h.DeleteEggPriceDiscount)
}

func NewEggPriceHandler(log *zap.Logger, service service.IEggPriceService, validator *validator.Validate) *EggPriceHandler {
	return &EggPriceHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *EggPriceHandler) CreateEggPrice(c *fiber.Ctx) error {
	var request dto.CreateEggPriceRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[CreateEggPrice] failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[CreateEggPrice] failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[CreateEggPrice] userId not found in locals")
		return errx.Unauthorized("userId not found in locals")
	}

	resp, err := h.service.CreateEggPrice(request, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[CreateEggPrice] failed to create egg price", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusCreated,
		resp,
		"success create egg price",
	)
}

func (h *EggPriceHandler) GetEggPrices(c *fiber.Ctx) error {
	eggPrices, err := h.service.GetEggPrices()
	if err != nil {
		h.log.Error("[GetEggPrices] failed to get egg prices", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		eggPrices,
		"success get egg prices",
	)
}

func (h *EggPriceHandler) GetEggPriceById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("[GetEggPriceById] failed to parse id", zap.Error(err))
		return err
	}

	resp, err := h.service.GetEggPriceById(id)
	if err != nil {
		h.log.Error("[GetEggPriceById] failed to get egg price by id", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		resp,
		"success get egg price by id",
	)
}

func (h *EggPriceHandler) UpdateEggPrice(c *fiber.Ctx) error {
	var request dto.UpdateEggPriceRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[UpdateEggPrice] failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[UpdateEggPrice] failed to validate request", zap.Error(err))
		return err
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("[UpdateEggPrice] failed to parse id", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[UpdateEggPrice] userId not found in locals")
		return errx.Unauthorized("userId not found in locals")
	}

	resp, err := h.service.UpdateEggPrice(id, request, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[UpdateEggPrice] failed to update egg price", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		resp,
		"success update egg price",
	)
}

func (h *EggPriceHandler) DeleteEggPrice(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("[DeleteEggPrice] failed to parse id", zap.Error(err))
		return err
	}

	err = h.service.DeleteEggPrice(id)
	if err != nil {
		h.log.Error("[DeleteEggPrice] failed to delete egg price", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}

func (h *EggPriceHandler) CreateEggPriceDiscount(c *fiber.Ctx) error {
	var filter dto.CreateEggPriceDiscountRequest
	if err := c.BodyParser(&filter); err != nil {
		h.log.Error("[CreateEggPriceDiscount] failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("[CreateEggPriceDiscount] failed to validate request", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[CreateEggPriceDiscount] userId not found in locals")
		return errx.Unauthorized("userId not found in locals")
	}

	resp, err := h.service.CreateEggPriceDiscount(filter, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[CreateEggPriceDiscount] failed to create egg price discount", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusCreated,
		resp,
		"success create egg price discount",
	)
}

func (h *EggPriceHandler) GetEggPriceDiscountById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("[GetEggPriceDiscountById] failed to parse id", zap.Error(err))
		return err
	}

	resp, err := h.service.GetEggPriceDiscountById(id)
	if err != nil {
		h.log.Error("[GetEggPriceDiscountById] failed to get egg price discount by id", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		resp,
		"success get egg price discount by id",
	)
}

func (h *EggPriceHandler) GetEggPriceDiscounts(c *fiber.Ctx) error {
	resp, err := h.service.GetEggPriceDiscounts()
	if err != nil {
		h.log.Error("[GetEggPriceDiscounts] failed to get egg price discounts", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		resp,
		"success get egg price discounts",
	)
}

func (h *EggPriceHandler) UpdateEggPriceDiscount(c *fiber.Ctx) error {
	var request dto.UpdateEggPriceDiscountRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("[UpdatePriceDiscount] failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("[UpdatePriceDiscount] failed to validate request", zap.Error(err))
		return err
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("[UpdatePriceDiscount] failed to parse id", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Error("[UpdatePriceDiscount] userId not found in locals")
		return errx.Unauthorized("userId not found in locals")
	}

	resp, err := h.service.UpdateEggPriceDiscount(id, request, uuid.MustParse(userId))
	if err != nil {
		h.log.Error("[UpdatePriceDiscount] failed to update egg price discount", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		resp,
		"success update egg price discount",
	)
}

func (h *EggPriceHandler) DeleteEggPriceDiscount(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("[DeleteEggPriceDiscount] failed to parse id", zap.Error(err))
		return err
	}

	err = h.service.DeleteEggPriceDiscount(id)
	if err != nil {
		h.log.Error("[DeleteEggPriceDiscount] failed to delete egg price discount", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}
