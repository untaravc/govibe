package homecontroller

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type HomeController struct{}

func New() *HomeController {
	return &HomeController{}
}

func (h *HomeController) Index(c *fiber.Ctx) error {
	fullWidth := strings.HasPrefix(c.Path(), "/admin")
	return c.Render("index", fiber.Map{
		"Title":     "GoVibe (Fiber)",
		"FullWidth": fullWidth,
	}, "layout")
}
