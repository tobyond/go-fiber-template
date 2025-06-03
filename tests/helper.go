package tests

import (
	"fresh/app/controllers"
	"fresh/app/models"
	"fresh/app/services"
	"fresh/routes"
	"log"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestApp holds the test application setup
type TestApp struct {
	App             *fiber.App
	DB              *gorm.DB
	UserRepo        *models.UserRepository
	AuthService     *services.AuthService
	TemplateService services.TemplateRenderer
	AuthController  *controllers.AuthController
	DashController  *controllers.DashboardController
}

// SetupTestApp creates a test application with in-memory SQLite database
func SetupTestApp(t *testing.T) *TestApp {
	// Create in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Quiet during tests
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Create mock template service for tests
	templateService := &MockTemplateService{}

	// Initialize services and controllers
	userRepo := models.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService, templateService)
	dashController := controllers.NewDashboardController(authService, templateService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		// Disable startup message in tests
		DisableStartupMessage: true,
		// Custom error handler for tests
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			t.Logf("Test app error: %v", err)
			return c.Status(500).SendString("Internal Server Error")
		},
	})

	// Setup routes
	routes.SetupRoutes(app, authController, dashController, authService)

	return &TestApp{
		App:             app,
		DB:              db,
		UserRepo:        userRepo,
		AuthService:     authService,
		TemplateService: templateService,
		AuthController:  authController,
		DashController:  dashController,
	}
}

// TeardownTestApp cleans up after tests
func (ta *TestApp) TeardownTestApp() {
	sqlDB, _ := ta.DB.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

// CreateTestUser creates a user for testing
func (ta *TestApp) CreateTestUser(email, password string) (*models.User, error) {
	return ta.UserRepo.Create(email, password)
}

// MockTemplateService is a mock template service for tests
type MockTemplateService struct{}

func (m *MockTemplateService) Render(c *fiber.Ctx, templateName string, data interface{}) error {
	// Mock template rendering for tests
	return c.JSON(fiber.Map{
		"template": templateName,
		"data":     data,
	})
}

// TestMain sets up and tears down for all tests in the package
func TestMain(m *testing.M) {
	// Setup
	log.SetOutput(os.Stdout) // Ensure logs are visible during tests

	// Run tests
	code := m.Run()

	// Cleanup and exit
	os.Exit(code)
}
