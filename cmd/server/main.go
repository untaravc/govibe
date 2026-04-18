package main

import (
	"log"

	accessmiddleware "govibe/app/Http/Middleware/AccessMiddleware"
	"govibe/app/Http/Response"
	"govibe/configs"
	"govibe/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

func main() {
	views := html.New("./resources/views", ".html")

	app := fiber.New(fiber.Config{
		AppName:      "govibe",
		Views:        views,
		ErrorHandler: response.ErrorHandler(),
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(accessmiddleware.New())

	app.Static("/static", "./public")

	db, err := configs.OpenGormMySQL()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	routes.RegisterWeb(app)
	routes.RegisterAPI(app, db)

	appCfg := configs.LoadAppConfig()

	log.Printf("listening on http://localhost:%s", appCfg.Port)
	if err := app.Listen(":" + appCfg.Port); err != nil {
		log.Fatal(err)
	}
}
