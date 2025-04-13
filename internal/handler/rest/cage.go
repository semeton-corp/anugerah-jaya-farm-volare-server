package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/service"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/response"
	"go.uber.org/zap"
)

type CageHandler struct {
	log       *zap.Logger
	validator *validator.Validate
	service   service.ICageService
}

func (a *CageHandler) SetEndpoint(router *fiber.App) {
	v1 := router.Group("api/v1/cages")
	v1.Get("/", a.GetCages)
}

func NewCageHandler(log *zap.Logger, service service.ICageService, validator *validator.Validate) *CageHandler {
	return &CageHandler{
		log:       log,
		service:   service,
		validator: validator,
	}
}

func (a *CageHandler) GetCages(c *fiber.Ctx) error {
	res, err := a.service.GetCages()
	if err != nil {
		a.log.Error("[GetCages] failed to get cages", zap.Error(err))
		return err
	}

	return response.SuccessResponse(
		c,
		fiber.StatusOK,
		res,
		"success get cages",
	)
}
