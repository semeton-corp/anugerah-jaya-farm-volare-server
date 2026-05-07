package router

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	recoverer "github.com/gofiber/fiber/v2/middleware/recover"

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

	router.Use(recoverer.New(recoverer.Config{
		EnableStackTrace: true,
	}))

	return router
}

func GlobalErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		if je, ok := err.(*json.UnmarshalTypeError); ok {
			return response.ErrorResponse(
				c,
				fiber.StatusBadRequest,
				je.Error(),
				"failed to parse json",
			)
		}

		if fe, ok := err.(*fiber.Error); ok {
			return response.ErrorResponse(
				c,
				fe.Code,
				fe.Message,
				fe.Error(),
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

		if ce, ok := err.(*errx.Errx); ok {
			return response.ErrorResponse(
				c,
				ce.Err.Code,
				ce.Err.Error(),
				ce.Error(),
			)
		}

		return response.ErrorResponse(
			c,
			fiber.ErrInternalServerError.Code,
			nil,
			fiber.ErrInternalServerError.Message,
		)
	}
}
