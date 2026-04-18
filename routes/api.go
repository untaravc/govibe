package routes

import (
	authcontroller "govibe/app/Http/Controllers/AuthController"
	usercontroller "govibe/app/Http/Controllers/UserController"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterAPI(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	authController := authcontroller.New(db)
	api.Post("/register", authController.Register)
	api.Post("/login", authController.Login)
	api.Get("/profile", authController.Profile)

	userController := usercontroller.New(db)
	users := api.Group("/users")
	users.Get("/", userController.Index)
	users.Get("/:id", userController.Show)
	users.Post("/", userController.Store)
	users.Put("/:id", userController.Update)
	users.Delete("/:id", userController.Destroy)
}
