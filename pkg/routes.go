package pkg

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes initializes routes for the application
func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("views/index", fiber.Map{
			"Title": "Hello, World!",
		})
	})
	// Public routes
	app.Post("/login", LoginHandler)
	app.Post("/register", RegisterHandler)
	app.Get("/profile", UserHandler)
	app.Post("/submit", func(c *fiber.Ctx) error {
		// Handle the htmx form submission here
		return c.SendString("Form submitted successfully!")
	})
	// Authenticated routes for admin
	authAdminGroup := app.Group("/dashboard/admin", AuthMiddleware("admin"))
	authAdminGroup.Get("/", AdminDashboardHandler)

	// Authenticated routes for user
	authUserGroup := app.Group("/dashboard/user", AuthMiddleware("user"))
	authUserGroup.Get("/", UserDashboardHandler)
}

// AdminDashboardHandler handles the admin dashboard
func AdminDashboardHandler(c *fiber.Ctx) error {
	userData := c.Locals("userData").(fiber.Map)
	// Render the HTML template with user data
	return c.Render("views/home/home-admin", fiber.Map{
		"user":    userData,
		"message": "Welcome to the admin panel",
	})
}

// UserDashboardHandler handles the user dashboard
func UserDashboardHandler(c *fiber.Ctx) error {
	userData := c.Locals("userData").(fiber.Map)

	// Render the HTML template with user data
	return c.Render("views/home/home-user", fiber.Map{
		"user":    userData,
		"message": "Welcome to the admin panel",
	})
}
