package routes

import (
	authcontroller "govibe/app/Http/Controllers/AuthController"
	categorycontroller "govibe/app/Http/Controllers/CategoryController"
	mediacontroller "govibe/app/Http/Controllers/MediaController"
	menucontroller "govibe/app/Http/Controllers/MenuController"
	menurolecontroller "govibe/app/Http/Controllers/MenuRoleController"
	officecontroller "govibe/app/Http/Controllers/OfficeController"
	postcontroller "govibe/app/Http/Controllers/PostController"
	rolecontroller "govibe/app/Http/Controllers/RoleController"
	usercontroller "govibe/app/Http/Controllers/UserController"
	"govibe/app/Http/Response"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterAPI(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api")
	api.Get("/health", func(c *fiber.Ctx) error {
		return response.OK(c, "ok", fiber.Map{"status": "ok"})
	})

	authController := authcontroller.New(db)
	api.Post("/register", authController.Register)
	api.Post("/login", authController.Login)
	api.Post("/logout", authController.Logout)
	api.Post("/request-reset-password", authController.RequestResetPassword)
	api.Get("/profile", authController.Profile)

	menuController := menucontroller.New(db)
	api.Get("/menu", menuController.Index)
	api.Get("/menus", menuController.List)

	mediaController := mediacontroller.New()
	api.Post("/upload", mediaController.Upload)

	userController := usercontroller.New(db)
	users := api.Group("/users")
	users.Get("/", userController.Index)
	users.Get("/:id", userController.Show)
	users.Post("/", userController.Store)
	users.Put("/:id", userController.Update)
	users.Delete("/:id", userController.Destroy)

	roleController := rolecontroller.New(db)
	roles := api.Group("/roles")
	roles.Get("/", roleController.Index)
	roles.Get("/:id", roleController.Show)
	roles.Post("/", roleController.Store)
	roles.Put("/:id", roleController.Update)
	roles.Delete("/:id", roleController.Destroy)

	postController := postcontroller.New(db)
	posts := api.Group("/posts")
	posts.Get("/", postController.Index)
	posts.Get("/:id", postController.Show)
	posts.Post("/", postController.Store)
	posts.Put("/:id", postController.Update)
	posts.Delete("/:id", postController.Destroy)

	categoryController := categorycontroller.New(db)
	categories := api.Group("/categories")
	categories.Get("/", categoryController.Index)
	categories.Get("/:id", categoryController.Show)
	categories.Post("/", categoryController.Store)
	categories.Put("/:id", categoryController.Update)
	categories.Delete("/:id", categoryController.Destroy)

	officeController := officecontroller.New(db)
	offices := api.Group("/offices")
	offices.Get("/", officeController.Index)
	offices.Get("/:id", officeController.Show)
	offices.Post("/", officeController.Store)
	offices.Put("/:id", officeController.Update)
	offices.Delete("/:id", officeController.Destroy)

	menuRoleController := menurolecontroller.New(db)
	api.Get("/menu-roles", menuRoleController.Index)
	api.Post("/menu-roles", menuRoleController.Save)
	api.Put("/menu-roles", menuRoleController.Save)
}
