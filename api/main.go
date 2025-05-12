package main

import (
	"log"
	"os"
	"api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupRoutes(app *fiber.App) {
	app.Post("/api/v1", routes.ShortenURL)
	app.Get("/api/v1/get/:url", routes.GetURL)
}	

func main() {

	app := fiber.New()

	app.Use(logger.New())
	setupRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}