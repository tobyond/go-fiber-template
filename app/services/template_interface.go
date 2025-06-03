package services

import "github.com/gofiber/fiber/v2"

// TemplateRenderer interface that both real and mock template services can implement
type TemplateRenderer interface {
	Render(c *fiber.Ctx, templateName string, data interface{}) error
}
