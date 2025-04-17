package middleware

import (
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
)

func Authorization(roles ...string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Get("role")
		if !slices.Contains(roles, c.Get("role")) {
			return errx.Unauthorized("invalid role")
		}

		return c.Next()
	}
}
