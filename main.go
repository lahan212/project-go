package main

import (
	"project-go/pkg"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars/v2"
)

func main() {

	// Register HTML template engine
	engine := handlebars.New("./web", ".hbs")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/static", "./web/styles")
	// Initialize routes
	pkg.SetupRoutes(app)

	// Start the Fiber app on port 3000
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
