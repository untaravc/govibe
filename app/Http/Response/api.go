package response

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Result  any    `json:"result"`
}

func JSON(c *fiber.Ctx, status int, success bool, message string, result any) error {
	if result == nil {
		result = fiber.Map{}
	}
	return c.Status(status).JSON(APIResponse{
		Success: success,
		Message: message,
		Result:  result,
	})
}

func OK(c *fiber.Ctx, message string, result any) error {
	return JSON(c, fiber.StatusOK, true, message, result)
}

func Created(c *fiber.Ctx, message string, result any) error {
	return JSON(c, fiber.StatusCreated, true, message, result)
}

func Error(c *fiber.Ctx, status int, message string, result any) error {
	return JSON(c, status, false, message, result)
}

