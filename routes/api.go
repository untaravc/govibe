package routes

import (
	authcontroller "govibe/app/Http/Controllers/AuthController"
	categorycontroller "govibe/app/Http/Controllers/CategoryController"
	mediacontroller "govibe/app/Http/Controllers/MediaController"
	menucontroller "govibe/app/Http/Controllers/MenuController"
	menurolecontroller "govibe/app/Http/Controllers/MenuRoleController"
	officecontroller "govibe/app/Http/Controllers/OfficeController"
	postcontroller "govibe/app/Http/Controllers/PostController"
	regioncontroller "govibe/app/Http/Controllers/RegionController"
	rolecontroller "govibe/app/Http/Controllers/RoleController"
	sectioncontroller "govibe/app/Http/Controllers/SectionController"
	shipmentcontroller "govibe/app/Http/Controllers/ShipmentController"
	usercontroller "govibe/app/Http/Controllers/UserController"
	authmiddleware "govibe/app/Http/Middleware/AuthMiddleware"
	"govibe/app/Http/Response"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterAPI(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api")
	api.Get("/health", func(c *fiber.Ctx) error { return response.OK(c, "ok", fiber.Map{"status": "ok"}) })

	authController := authcontroller.New(db)
	// Public auth endpoints (no bearer token required).
	api.Post("/register", authController.Register)
	api.Post("/login", authController.Login)
	api.Post("/refresh-token", authController.RefreshToken)
	api.Post("/request-reset-password", authController.RequestResetPassword)
	api.Get("/validate-email-token", authController.ValidateEmailToken)
	api.Post("/update-password-with-token", authController.UpdatePasswordWithToken)

	// Protected API endpoints (bearer token required).
	protected := api.Group("", authmiddleware.New(db))
	protected.Post("/logout", authController.Logout)
	protected.Get("/profile", authController.Profile)
	protected.Patch("/profile", authController.UpdateProfile)

	menuController := menucontroller.New(db)
	protected.Get("/menu", menuController.Index)
	protected.Get("/menus", menuController.List)

	mediaController := mediacontroller.New()
	protected.Post("/upload", mediaController.Upload)

	sectionController := sectioncontroller.New()
	protected.Get("/sections", sectionController.Index)

	regionController := regioncontroller.New(db)
	protected.Get("/provinces", regionController.Provinces)
	protected.Get("/cities", regionController.Cities)
	protected.Get("/districts", regionController.Districts)
	protected.Get("/villages", regionController.Villages)

	userController := usercontroller.New(db)
	users := protected.Group("/users")
	users.Get("/", userController.Index)
	users.Get("/:id", userController.Show)
	users.Post("/", userController.Store)
	users.Put("/:id", userController.Update)
	users.Delete("/:id", userController.Destroy)

	roleController := rolecontroller.New(db)
	roles := protected.Group("/roles")
	roles.Get("/", roleController.Index)
	roles.Get("/:id", roleController.Show)
	roles.Post("/", roleController.Store)
	roles.Put("/:id", roleController.Update)
	roles.Delete("/:id", roleController.Destroy)

	postController := postcontroller.New(db)
	posts := protected.Group("/posts")
	posts.Get("/", postController.Index)
	posts.Get("/:id", postController.Show)
	posts.Post("/", postController.Store)
	posts.Put("/:id", postController.Update)
	posts.Delete("/:id", postController.Destroy)

	categoryController := categorycontroller.New(db)
	categories := protected.Group("/categories")
	categories.Get("/", categoryController.Index)
	categories.Get("/:id", categoryController.Show)
	categories.Post("/", categoryController.Store)
	categories.Put("/:id", categoryController.Update)
	categories.Delete("/:id", categoryController.Destroy)

	shipmentController := shipmentcontroller.New(db)
	shipments := protected.Group("/shipments")
	shipments.Get("/", shipmentController.Index)
	shipments.Get("/:id", shipmentController.Show)
	shipments.Post("/", shipmentController.Store)
	shipments.Put("/:id", shipmentController.Update)
	shipments.Delete("/:id", shipmentController.Destroy)

	officeController := officecontroller.New(db)
	offices := protected.Group("/offices")
	offices.Get("/", officeController.Index)
	offices.Get("/:id", officeController.Show)
	offices.Post("/", officeController.Store)
	offices.Put("/:id", officeController.Update)
	offices.Delete("/:id", officeController.Destroy)

	menuRoleController := menurolecontroller.New(db)
	protected.Get("/menu-roles", menuRoleController.Index)
	protected.Post("/menu-roles", menuRoleController.Save)
	protected.Put("/menu-roles", menuRoleController.Save)
}
