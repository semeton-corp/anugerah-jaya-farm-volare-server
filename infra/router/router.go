package router

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
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
			Prefork:       true,
			CaseSensitive: true,
		},
	)

	router.Use(cors.New(cors.Config{
		AllowMethods:  viper.GetString("server.cors.allow_methods"),
		AllowHeaders:  viper.GetString("server.cors.allow_headers"),
		ExposeHeaders: viper.GetString("server.cors.expose_headers"),
		MaxAge:        viper.GetInt("server.cors.max_age"),
	}))

	return router
}

func GlobalErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		if ce, ok := err.(*errx.Errx); ok {
			return c.Status(ce.Err.Code).JSON(fiber.Map{
				"message": ce.Err.Message,
				"error":   ce.Error(),
				"status":  ce.Err.Code,
			})
		}

		if ve, ok := err.(validator.ValidationErrors); ok {
			out := make(map[string]string)
			for _, e := range ve {
				out[e.Field()] = util.GetErrorValidationMessage(e)
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": fiber.ErrBadRequest.Message,
				"error":   out,
				"status":  fiber.StatusBadRequest,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fiber.ErrInternalServerError.Message,
			"error":   err.Error(),
			"status":  fiber.StatusInternalServerError,
		})
	}
}
