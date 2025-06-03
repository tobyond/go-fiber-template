package controllers

import (
	"fresh/app/services"

	"github.com/gofiber/fiber/v2"
)

type DashboardController struct {
	authService     *services.AuthService
	templateService services.TemplateRenderer
}

func NewDashboardController(authService *services.AuthService, templateService services.TemplateRenderer) *DashboardController {
	return &DashboardController{
		authService:     authService,
		templateService: templateService,
	}
}

func (dc *DashboardController) Show(c *fiber.Ctx) error {
	// Check authentication
	user, err := dc.authService.GetCurrentUser(c)
	if err != nil {
		return c.Redirect("/login")
	}

	return dc.templateService.Render(c, "dashboard", fiber.Map{
		"Title": "Dashboard - Fresh",
		"User":  user,
	})
}
