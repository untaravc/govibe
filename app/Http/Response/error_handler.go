package response

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		if strings.HasPrefix(c.Path(), "/api") {
			status := fiber.StatusInternalServerError
			message := "internal server error"

			if e, ok := err.(*fiber.Error); ok && e != nil {
				status = e.Code
				message = e.Message
			} else if err != nil {
				message = err.Error()
			}

			return JSON(c, status, false, message, fiber.Map{})
		}

		return fiber.DefaultErrorHandler(c, err)
	}
}

