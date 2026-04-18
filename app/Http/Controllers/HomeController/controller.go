package homecontroller

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type HomeController struct{}

func New() *HomeController {
	return &HomeController{}
}

func (h *HomeController) Index(c *fiber.Ctx) error {
	fullWidth := strings.HasPrefix(c.Path(), "/admin")
	return c.Render("index", fiber.Map{
		"Title":        "GoVibe (Fiber)",
		"FullWidth":    fullWidth,
		"AssetVersion": getAssetVersion(),
	}, "layout")
}

func getAssetVersion() string {
	if fi, err := os.Stat("public/dist/app.js"); err == nil {
		return fmt.Sprintf("%d", fi.ModTime().UTC().Unix())
	}
	if fi, err := os.Stat("./public/dist/app.js"); err == nil {
		return fmt.Sprintf("%d", fi.ModTime().UTC().Unix())
	}
	return fmt.Sprintf("%d", time.Now().UTC().Unix())
}
