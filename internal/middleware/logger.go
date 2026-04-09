package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"go.uber.org/zap"
)

func RequestLogger(log *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()
		if err != nil {
			switch e := err.(type) {
			case *fiber.Error:
				status = e.Code
			case *errx.Errx:
				status = e.Err.Code
			default:
				if status < fiber.StatusBadRequest {
					status = fiber.StatusInternalServerError
				}
			}
		}

		fields := []zap.Field{
			zap.String("ip", c.IP()),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("user_agent", c.Get("User-Agent")),
			zap.Time("time", time.Now()),
		}

		if err != nil {
			fields = append(fields, zap.Error(err))
		}

		log.Info("request", fields...)

		return err
	}
}
