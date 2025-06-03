package routes

import (
	"fresh/app/controllers"
	"fresh/app/middleware"
	"fresh/app/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authController *controllers.AuthController, dashboardController *controllers.DashboardController, authService *services.AuthService) {
	// Root redirect
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/login")
	})

	// Public routes (no middleware for now to debug)
	app.Get("/login", authController.ShowLogin)
	app.Post("/login", authController.HandleLogin)
	app.Get("/register", authController.ShowRegister)
	app.Post("/register", authController.HandleRegister)

	// Protected routes (require authentication)
	protected := app.Group("/", middleware.RequireAuth(authService))
	protected.Get("/dashboard", dashboardController.Show)

	// Logout (no middleware needed)
	app.Post("/logout", authController.HandleLogout)
}
