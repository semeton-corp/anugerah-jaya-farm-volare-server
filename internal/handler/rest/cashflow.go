package rest

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/service"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/response"
	"go.uber.org/zap"
)

type CashflowHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.ICashflowService
}

func (h *CashflowHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("/cashflow")

	v1.Get("/income/overview", h.GetIncomeOverview)
	v1.Get("/income/:category/:id", h.GetIncome)

}

func NewCashflowHandler(log *zap.Logger, service service.ICashflowService, validator *validator.Validate) *CashflowHandler {
	return &CashflowHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (h *CashflowHandler) GetIncomeOverview(c *fiber.Ctx) error {
	var filter dto.GetIncomeOverviewFilter

	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query params", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid query parameters",
		})
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Warn("validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resp, err := h.service.GetIncomeOverview(filter)
	if err != nil {
		h.log.Error("service error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get income overview",
		})
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get income overview")
}

func (h *CashflowHandler) GetIncome(c *fiber.Ctx) error {
	category := c.Params("category")
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	resp, err := h.service.GetIncome(category, id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get income")
}

func (h *CashflowHandler) ExportSalesToExcel(c *fiber.Ctx) error {
	var filter dto.GetSaleCashflowFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed parse query", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&filter); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return err
	}

	f, err := h.service.ExportSalesCashflowToExcel(filter)
	if err != nil {
		h.log.Error("failed to export sales", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to export sales",
		})
	}

	fileName := fmt.Sprintf("sales_report_%s.xlsx", time.Now().Format("20060102_150405"))

	buf, err := f.WriteToBuffer()
	if err != nil {
		h.log.Error("failed to write excel buffer", zap.Error(err))
		return err
	}

	// set headers for file download
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	return c.SendStream(buf)
}
