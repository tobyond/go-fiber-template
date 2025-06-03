package services

import (
	"fmt"
	"html/template"

	"github.com/gofiber/fiber/v2"
)

type TemplateService struct {
	templates map[string]*template.Template
}

func NewTemplateService() (*TemplateService, error) {
	ts := &TemplateService{
		templates: make(map[string]*template.Template),
	}

	err := ts.parsePageTemplates()
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func (ts *TemplateService) parsePageTemplates() error {
	// Create separate template instances for each page to avoid conflicts
	pages := []string{"login", "register", "dashboard"}

	for _, page := range pages {
		// Create a new template instance for this page
		t := template.New(page)

		// Parse layout first
		if _, err := t.ParseFiles("web/templates/layout.html"); err != nil {
			return fmt.Errorf("parsing layout: %v", err)
		}

		// Parse the specific page template
		pageFile := fmt.Sprintf("web/templates/%s.html", page)
		if _, err := t.ParseFiles(pageFile); err != nil {
			return fmt.Errorf("parsing %s: %v", pageFile, err)
		}

		// Store this template instance
		ts.templates[page] = t
	}

	return nil
}

func (ts *TemplateService) Render(c *fiber.Ctx, templateName string, data interface{}) error {
	t, exists := ts.templates[templateName]
	if !exists {
		return fmt.Errorf("template %s not found", templateName)
	}

	c.Set("Content-Type", "text/html")

	// Execute the layout template (which will include the page content)
	return t.ExecuteTemplate(c.Response().BodyWriter(), "layout", data)
}
