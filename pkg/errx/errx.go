package errx

import "github.com/gofiber/fiber/v2"

type Errx struct {
	Err     *fiber.Error
	Message string
}

func (e *Errx) Error() string {
	return e.Message
}

func BadRequest(s string) error {
	return &Errx{
		Err:     fiber.ErrBadRequest,
		Message: s,
	}
}

func NotFound(s string) error {
	return &Errx{
		Err:     fiber.ErrNotFound,
		Message: s,
	}
}

func InternalServerError(s string) error {
	return &Errx{
		Err:     fiber.ErrInternalServerError,
		Message: s,
	}
}

func Unauthorized(s string) error {
	return &Errx{
		Err:     fiber.ErrUnauthorized,
		Message: s,
	}
}

func Forbidden(s string) error {
	return &Errx{
		Err:     fiber.ErrForbidden,
		Message: s,
	}
}
