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

func (a *SupplierHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/suppliers")

	v1.Post("/", middleware.Authentication(), a.CreateSupplier)
	v1.Get("/:id", middleware.Authentication(), a.GetSupplierById)
	v1.Get("/", middleware.Authentication(), a.GetAllSuppliers)
	v1.Get("", middleware.Authentication(), a.GetAllSuppliers)
	v1.Put("/:id", middleware.Authentication(), a.UpdateSupplier)
	v1.Delete("/:id", middleware.Authentication(), a.DeleteSupplier)
}

func NewSupplierHandler(log *zap.Logger, service service.ISupplierService, validator *validator.Validate) *SupplierHandler {
	return &SupplierHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (a *SupplierHandler) CreateSupplier(c *fiber.Ctx) error {
	var request dto.CreateSupplierRequest
	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[CreateSupplier] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[CreateSupplier] failed to validate request", zap.Error(err))
		return err
	}

	accountIdCtx := c.Locals("accountId").(string)
	if accountIdCtx == "" {
		a.log.Error("[CreateSupplier] accountId not found in context")
		return errx.NotFound("accountId not found in context")
	}

	resp, err := a.service.CreateSupplier(&request, uuid.MustParse(accountIdCtx))
	if err != nil {
		a.log.Error("[CreateSupplier] failed to create supplier", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, resp, "success create supplier")
}

func (a *SupplierHandler) GetSupplierById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		a.log.Error("[GetSupplierById] failed to parse id", zap.Error(err))
		return err
	}

	supplier, err := a.service.GetSupplierById(id)
	if err != nil {
		a.log.Error("[GetSupplierById] failed to get supplier", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, supplier, "success get supplier by id")
}

func (a *SupplierHandler) GetAllSuppliers(c *fiber.Ctx) error {
	suppliers, err := a.service.GetAllSuppliers()
	if err != nil {
		a.log.Error("[GetAllSuppliers] failed to get all suppliers", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, suppliers, "success get all suppliers")
}

func (a *SupplierHandler) UpdateSupplier(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		a.log.Error("[UpdateSupplier] failed to parse id", zap.Error(err))
		return err
	}

	accountIdCtx := c.Locals("accountId").(string)
	if accountIdCtx == "" {
		a.log.Error("[UpdateSupplier] accountId not found in context")
		return errx.NotFound("accountId not found in context")
	}

	var request dto.UpdateSupplierRequest

	if err := c.BodyParser(&request); err != nil {
		a.log.Error("[UpdateSupplier] failed to parse request", zap.Error(err))
		return err
	}

	if err := a.validator.Struct(request); err != nil {
		a.log.Error("[UpdateSupplier] failed to validate request", zap.Error(err))
		return err
	}

	resp, err := a.service.UpdateSupplier(id, &request, uuid.MustParse(accountIdCtx))
	if err != nil {
		a.log.Error("[UpdateSupplier] failed to update supplier", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success update supplier")
}

func (a *SupplierHandler) DeleteSupplier(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		a.log.Error("[DeleteSupplier] failed to parse id", zap.Error(err))
		return err
	}

	if err := a.service.DeleteSupplier(id); err != nil {
		a.log.Error("[DeleteSupplier] failed to delete supplier", zap.Error(err))
		return err
	}

	return response.NoContentResponse(c)
}
