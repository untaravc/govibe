package routes

import (
	homecontroller "govibe/app/Http/Controllers/HomeController"

	"github.com/gofiber/fiber/v2"
)

func RegisterWeb(app *fiber.App) {
	homeController := homecontroller.New()
	app.Get("/", homeController.Index)
	// SPA fallback for Vue Router history mode (only for frontend prefixes).
	app.Get("/admin", homeController.Index)
	app.Get("/admin/*", homeController.Index)
	app.Get("/auth", homeController.Index)
	app.Get("/auth/*", homeController.Index)
}
