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

type SupplierHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.ISupplierService
}

func (h *SupplierHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/suppliers")

	v1.Post("/", middleware.Authentication(), h.CreateSupplier)
	v1.Get("/:id", middleware.Authentication(), h.GetSupplierById)
	v1.Get("/", middleware.Authentication(), h.GetAllSuppliers)
	v1.Get("", middleware.Authentication(), h.GetAllSuppliers)
	v1.Put("/:id", middleware.Authentication(), h.UpdateSupplier)
	v1.Delete("/:id", middleware.Authentication(), h.DeleteSupplier)
}

func NewSupplierHandler(log *zap.Logger, service service.ISupplierService, validator *validator.Validate) *SupplierHandler {
	return &SupplierHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *SupplierHandler) CreateSupplier(c *fiber.Ctx) error {
	var request dto.CreateSupplierRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	userId := c.Locals("userId").(string)
	if userId == "" {
		h.log.Error("user id not found in context")
		return errx.NotFound("user id not found in context")
	}

	resp, err := h.service.CreateSupplier(&request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, resp, "success create supplier")
}

func (h *SupplierHandler) GetSupplierById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return err
	}

	supplier, err := h.service.GetSupplierById(id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, supplier, "success get supplier by id")
}

func (h *SupplierHandler) GetAllSuppliers(c *fiber.Ctx) error {
	suppliers, err := h.service.GetAllSuppliers()
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, suppliers, "success get all suppliers")
}

func (h *SupplierHandler) UpdateSupplier(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return err
	}

	userId := c.Locals("userId").(string)
	if userId == "" {
		h.log.Error("user id not found in context")
		return errx.NotFound("user id not found in context")
	}

	var request dto.UpdateSupplierRequest

	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("failed to validate request", zap.Error(err))
		return err
	}

	resp, err := h.service.UpdateSupplier(id, &request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success update supplier")
}

func (h *SupplierHandler) DeleteSupplier(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("failed to parse id", zap.Error(err))
		return err
	}

	if err := h.service.DeleteSupplier(id); err != nil {
		h.log.Error("failed to delete supplier", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}
