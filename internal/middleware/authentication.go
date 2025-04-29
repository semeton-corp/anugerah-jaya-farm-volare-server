package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/jwt"
	"go.uber.org/zap"
)

func Authentication() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")

		if authorization == "" {
			zap.L().Warn("missing authorization header")
			return errx.Unauthorized("missing token")
		}

		authorizations := strings.SplitN(authorization, " ", 2)
		if len(authorizations) != 2 || authorizations[0] != "Bearer" {
			zap.L().Warn("invalid authorization header format")
			return errx.Unauthorized("invalid token format")
		}

		token := authorizations[1]
		payload, err := jwt.DecodeToken(token)
		if err != nil {
			zap.L().Error("failed to decode token", zap.Error(err))
			return err
		}

		c.Locals("accountId", payload.ID)

		return c.Next()
	}
}
