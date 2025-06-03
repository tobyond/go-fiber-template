package middleware

import (
	"fresh/app/services"

	"github.com/gofiber/fiber/v2"
)

func RequireAuth(authService *services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := authService.GetCurrentUser(c)
		if err != nil {
			return c.Redirect("/login")
		}
		return c.Next()
	}
}

func RedirectIfAuthenticated(authService *services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := authService.GetCurrentUser(c)
		if err == nil {
			// User is authenticated, redirect to dashboard
			return c.Redirect("/dashboard")
		}
		return c.Next()
	}
}
