package router

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/response"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
	"github.com/spf13/viper"
)

func New() *fiber.App {
	router := fiber.New(
		fiber.Config{
			WriteTimeout:  viper.GetDuration("server.write_timeout"),
			ReadTimeout:   viper.GetDuration("server.read_timeout"),
			AppName:       viper.GetString("app.name"),
			ErrorHandler:  GlobalErrorHandler(),
			Prefork:       false,
			CaseSensitive: true,
		},
	)

	router.Use(logger.New(logger.Config{
		Format:     "${time} ${status} ${latency} ${method} ${path}\n",
		TimeFormat: "02-01-2006 15:04:05",
		TimeZone:   "Asia/Jakarta",
	}))

	return router
}

func GlobalErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		if fe, ok := err.(*fiber.Error); ok {
			return response.ErrorResponse(
				c,
				fe.Code,
				fe.Message,
				fe.Error(),
			)
		}

		if ce, ok := err.(*errx.Errx); ok {
			return response.ErrorResponse(
				c,
				ce.Err.Code,
				ce.Err.Error(),
				ce.Error(),
			)
		}

		if ve, ok := err.(validator.ValidationErrors); ok {
			out := make(map[string]string)
			for _, e := range ve {
				out[e.Field()] = util.GetErrorValidationMessage(e)
			}

			return response.ErrorResponse(
				c,
				fiber.ErrBadRequest.Code,
				out,
				fiber.ErrBadRequest.Message,
			)
		}

		return response.ErrorResponse(
			c,
			fiber.ErrInternalServerError.Code,
			err.Error(),
			fiber.ErrInternalServerError.Message,
		)
	}
}
