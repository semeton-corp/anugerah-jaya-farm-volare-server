package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Status  int    `json:"status"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
	Message string `json:"message"`
}

func SuccessResponse(ctx *fiber.Ctx, status int, data any, message string) error {
	return ctx.Status(status).JSON(Response{
		Status:  status,
		Data:    data,
		Message: message,
	})
}

func NoContentResponse(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusNoContent)
}

func ErrorResponse(ctx *fiber.Ctx, status int, err any, message string) error {
	return ctx.Status(status).JSON(Response{
		Status:  status,
		Error:   err,
		Message: message,
	})
}
