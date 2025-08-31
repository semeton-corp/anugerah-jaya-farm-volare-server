package rest

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

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

type CashflowHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.ICashflowService
}

func (h *CashflowHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("/api/v1/cashflows")

	v1.Get("/sales/overview", middleware.Authentication(), h.GetCashflowSaleOverview)
	v1.Get("/overview", middleware.Authentication(), h.GetCashflowOverview)

	v1.Get("/incomes/:category/:id", middleware.Authentication(), h.GetIncome)
	v1.Get("/incomes/overview", middleware.Authentication(), h.GetIncomeOverview)

	v1.Post("/expenses", middleware.Authentication(), h.CreateExpense)
	v1.Get("/expenses/:category/:id", middleware.Authentication(), h.GetExpense)
	v1.Get("/expenses/overview", middleware.Authentication(), h.GetExpenseOverview)

	v1.Get("/cash-advances/:userId", middleware.Authentication(), h.GetUserCashAdvanceByUserId)
	v1.Post("/cash-advances", middleware.Authentication(), h.CreateUserCashAdvance)
	v1.Post("/cash-advances/:id/payments", middleware.Authentication(), h.CreateUserCashAdvancePayment)

	v1.Get("/receivables/overview", middleware.Authentication(), h.GetReceivablesOverview)
	v1.Get("/receivables/:category/:id", middleware.Authentication(), h.GetReceivables)

	v1.Get("/debts/overview", middleware.Authentication(), h.GetDebtOverview)
	v1.Get("/debts/:category/:id", middleware.Authentication(), h.GetDebt)

	v1.Post("/salaries/:id/pay", middleware.Authentication(), h.PayUserSalaryPayment)
	v1.Get("/salaries/summary", middleware.Authentication(), h.GetUserSalarySummary)
	v1.Get("/salaries", middleware.Authentication(), h.GetUserSalaries)
	v1.Get("/salaries/:id", middleware.Authentication(), h.GetUserSalaryDetail)

	v1.Get("/sales/reports", middleware.Authentication(), h.ExportSalesToExcel)
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
		h.log.Error("error validation", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resp, err := h.service.GetIncomeOverview(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get income overview")
}

func (h *CashflowHandler) GetIncome(c *fiber.Ctx) error {
	category, err := url.QueryUnescape(c.Params("category"))
	if err != nil {
		h.log.Error("failed to decode category param", zap.Error(err))
		return errx.BadRequest("invalid category param")
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
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
	var filter dto.GetCashflowSaleReportFilter
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
		return err
	}

	fileName := fmt.Sprintf("sales_report_%s.xlsx", time.Now().Format("20060102_150405"))

	buf, err := f.WriteToBuffer()
	if err != nil {
		h.log.Error("failed to write excel buffer", zap.Error(err))
		return err
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	return c.SendStream(buf)
}

func (h *CashflowHandler) CreateExpense(c *fiber.Ctx) error {
	var request dto.CreateExpenseRequest

	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed to parse body", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&request); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		h.log.Warn("user id not found in context")
		return errx.Unauthorized("user id not found in context")
	}

	resp, err := h.service.CreateExpense(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, resp, "success create expense")
}

func (h *CashflowHandler) GetExpenseOverview(c *fiber.Ctx) error {
	var filter dto.GetExpenseOverviewFilter

	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed to parse query params", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return err
	}

	resp, err := h.service.GetExpenseOverview(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get expense overview")
}

func (h *CashflowHandler) GetExpense(c *fiber.Ctx) error {
	category, err := url.QueryUnescape(c.Params("category"))
	if err != nil {
		h.log.Error("failed to decode category param", zap.Error(err))
		return errx.BadRequest("invalid category param")
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	resp, err := h.service.GetExpense(category, id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get expense")
}

func (h *CashflowHandler) GetUserCashAdvanceByUserId(c *fiber.Ctx) error {
	userId := c.Params("userId")
	resp, err := h.service.GetUserCashAdvanceByUserId(uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get user cash advances")
}

func (h *CashflowHandler) CreateUserCashAdvance(c *fiber.Ctx) error {
	var request dto.CreateUserCashAdvanceRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed parse body request", zap.Error(err))
		return errx.BadRequest("failed to parse request body")
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return errx.BadRequest(err.Error())
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return errx.Unauthorized("user id not found in context")
	}

	resp, err := h.service.CreateUserCashAdvance(request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, resp, "success create user cash advance")
}

func (h *CashflowHandler) CreateUserCashAdvancePayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	var request dto.CreateUserCashAdvancePaymentRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed parse body request", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return err
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return errx.Unauthorized("user id not found in context")
	}

	resp, err := h.service.CreateUserCashAdvancePayment(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusCreated, resp, "success create cash advance payment")
}

func (h *CashflowHandler) GetReceivablesOverview(c *fiber.Ctx) error {
	var filter dto.GetReceivablesOverviewFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed parse query param", zap.Error(err))
		return errx.BadRequest("invalid query parameters")
	}

	if err := h.validator.Struct(&filter); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return errx.BadRequest(err.Error())
	}

	resp, err := h.service.GetReceiveablesOverview(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get receivables overview")
}

func (h *CashflowHandler) GetReceivables(c *fiber.Ctx) error {
	category, err := url.QueryUnescape(c.Params("category"))
	if err != nil {
		h.log.Error("failed to decode category param", zap.Error(err))
		return errx.BadRequest("invalid category param")
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	resp, err := h.service.GetReceiveables(category, id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get receivables")
}

func (h *CashflowHandler) PayUserSalaryPayment(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	var request dto.PayUserSalaryPaymentRequest
	if err := c.BodyParser(&request); err != nil {
		h.log.Error("failed parse body request", zap.Error(err))
		return errx.BadRequest("failed to parse request body")
	}

	if err := h.validator.Struct(request); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return errx.BadRequest(err.Error())
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return errx.Unauthorized("user id not found in context")
	}

	resp, err := h.service.PayUserSalaryPayment(id, request, uuid.MustParse(userId))
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success pay user salary")
}

func (h *CashflowHandler) GetDebtOverview(c *fiber.Ctx) error {
	var filter dto.GetDebtOverviewFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed parse query param", zap.Error(err))
		return errx.BadRequest("invalid query parameters")
	}

	if err := h.validator.Struct(&filter); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return errx.BadRequest(err.Error())
	}

	resp, err := h.service.GetDebtOverview(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get debt overview")
}

func (h *CashflowHandler) GetDebt(c *fiber.Ctx) error {
	category, err := url.QueryUnescape(c.Params("category"))
	if err != nil {
		h.log.Error("failed to decode category param", zap.Error(err))
		return errx.BadRequest("invalid category param")
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	resp, err := h.service.GetDebt(category, id)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get debt")
}

func (h *CashflowHandler) GetUserSalarySummary(c *fiber.Ctx) error {
	var filter dto.GetUserSalarySummaryFilter

	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed parse query", zap.Error(err))
		return errx.BadRequest("invalid query params")
	}

	if err := h.validator.Struct(&filter); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return errx.BadRequest(err.Error())
	}

	resp, err := h.service.GetUserSalarySummary(filter)
	if err != nil {
		h.log.Error("service error", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get user salary summary")
}

func (h *CashflowHandler) GetUserSalaries(c *fiber.Ctx) error {
	var filter dto.GetUserSalaryListFilter

	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed parse query", zap.Error(err))
		return errx.BadRequest("invalid query params")
	}

	if err := h.validator.Struct(filter); err != nil {
		h.log.Error("validation error", zap.Error(err))
		return errx.BadRequest(err.Error())
	}

	resp, err := h.service.GetUserSalaries(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get user salaries")
}

func (h *CashflowHandler) GetUserSalaryDetail(c *fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.log.Error("invalid id param", zap.Error(err))
		return errx.BadRequest("invalid id param")
	}

	resp, err := h.service.GetUserSalaryDetail(id)
	if err != nil {
		h.log.Error("service error", zap.Error(err))
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, resp, "success get user salary detail")
}

func (h *CashflowHandler) GetCashflowSaleOverview(c *fiber.Ctx) error {
	var filter dto.GetCashflowSaleOverviewFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed parse query param", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&filter); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return err
	}

	data, err := h.service.GetCashflowSaleOverview(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get cashflow sale overview")
}

func (h *CashflowHandler) GetCashflowOverview(c *fiber.Ctx) error {
	var filter dto.GetCashflowOverviewFilter
	if err := c.QueryParser(&filter); err != nil {
		h.log.Error("failed parse query param", zap.Error(err))
		return err
	}

	if err := h.validator.Struct(&filter); err != nil {
		h.log.Error("error validation", zap.Error(err))
		return err
	}

	data, err := h.service.GetCashflowOverview(filter)
	if err != nil {
		return err
	}

	return response.SuccessResponse(c, fiber.StatusOK, data, "success get cashflow overview")
}
