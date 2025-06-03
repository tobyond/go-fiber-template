package main

import (
	"fmt"
	"fresh/app/controllers"
	"fresh/app/models"
	"fresh/app/services"
	"fresh/config"
	"fresh/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Initialize database
	db := config.InitDatabase()

	// Initialize template service
	templateService, err := services.NewTemplateService()
	if err != nil {
		log.Fatal("Failed to initialize templates:", err)
	}

	// Initialize repositories
	userRepo := models.NewUserRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService, templateService)
	dashboardController := controllers.NewDashboardController(authService, templateService)

	// Create new Fiber instance
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			fmt.Printf("Error: %v\n", err)
			return c.Status(code).SendString(err.Error())
		},
	})

	// Add logging middleware
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} ${query} | ${ip} | ${latency}\n",
		TimeFormat: "15:04:05",
		TimeZone:   "Local",
	}))

	// Static files
	app.Static("/static", "./web/static")

	// Setup routes
	routes.SetupRoutes(app, authController, dashboardController, authService)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Fresh server running on port %s\n", port)
	log.Fatal(app.Listen(":" + port))
}
