package controllers

import (
	"fmt"
	"fresh/app/services"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService     *services.AuthService
	templateService services.TemplateRenderer
}

func NewAuthController(authService *services.AuthService, templateService services.TemplateRenderer) *AuthController {
	return &AuthController{
		authService:     authService,
		templateService: templateService,
	}
}

func (ac *AuthController) ShowLogin(c *fiber.Ctx) error {
	return ac.templateService.Render(c, "login", fiber.Map{
		"Title": "Login - Fresh",
	})
}

func (ac *AuthController) HandleLogin(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Debug logging
	fmt.Printf("Login attempt - Email: %s, IP: %s\n", email, c.IP())

	user, err := ac.authService.Login(email, password)
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		return ac.templateService.Render(c, "login", fiber.Map{
			"Title": "Login - Fresh",
			"Error": err.Error(),
		})
	}

	ac.authService.SetUserSession(c, user)
	fmt.Printf("Login successful for user ID %d (email: %s)\n", user.ID, email)

	return c.Redirect("/dashboard")
}

func (ac *AuthController) ShowRegister(c *fiber.Ctx) error {
	return ac.templateService.Render(c, "register", fiber.Map{
		"Title": "Register - Fresh",
	})
}

func (ac *AuthController) HandleRegister(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	fmt.Printf("Registration attempt - Email: %s, IP: %s\n", email, c.IP())

	user, err := ac.authService.Register(email, password)
	if err != nil {
		fmt.Printf("Registration failed: %v\n", err)
		return ac.templateService.Render(c, "register", fiber.Map{
			"Title": "Register - Fresh",
			"Error": err.Error(),
		})
	}

	fmt.Printf("Registration successful for user ID %d (email: %s)\n", user.ID, email)

	// Automatically log in the user after successful registration
	ac.authService.SetUserSession(c, user)
	fmt.Printf("User automatically logged in after registration\n")

	return c.Redirect("/dashboard")
}

func (ac *AuthController) HandleLogout(c *fiber.Ctx) error {
	ac.authService.ClearUserSession(c)
	fmt.Printf("User logged out, IP: %s\n", c.IP())
	return c.Redirect("/login")
}
