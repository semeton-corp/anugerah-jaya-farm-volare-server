package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/jwt"
)

func Authentication(role string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")

		if authorization == "" {
			return errx.Unauthorized("missing token")
		}

		authorizations := strings.SplitN(authorization, " ", 2)
		if len(authorizations) != 2 || authorizations[0] != "Bearer" {
			return errx.Unauthorized("invalid token format")
		}

		token := authorizations[1]
		payload, err := jwt.DecodeToken(token)
		if err != nil {
			return err
		}

		if payload.Role != role {
			return errx.Forbidden("access denied")
		}

		c.Locals("accountId", payload.ID)

		return c.Next()
	}
}
