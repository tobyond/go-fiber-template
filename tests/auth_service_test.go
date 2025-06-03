package tests

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Register(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.TeardownTestApp()

	tests := []struct {
		name        string
		email       string
		password    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid registration",
			email:       "register@example.com",
			password:    "password123",
			expectError: false,
		},
		{
			name:        "empty email",
			email:       "",
			password:    "password123",
			expectError: true,
			errorMsg:    "email and password are required",
		},
		{
			name:        "empty password",
			email:       "test@example.com",
			password:    "",
			expectError: true,
			errorMsg:    "email and password are required",
		},
		{
			name:        "duplicate email",
			email:       "register@example.com", // Same as first test
			password:    "password456",
			expectError: true,
			errorMsg:    "email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := testApp.AuthService.Register(tt.email, tt.password)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
				assert.Greater(t, user.ID, uint(0))
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.TeardownTestApp()

	// Create a test user for login tests
	testEmail := "login@example.com"
	testPassword := "correctpassword"
	_, err := testApp.CreateTestUser(testEmail, testPassword)
	require.NoError(t, err)

	tests := []struct {
		name        string
		email       string
		password    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid login",
			email:       testEmail,
			password:    testPassword,
			expectError: false,
		},
		{
			name:        "invalid email",
			email:       "nonexistent@example.com",
			password:    testPassword,
			expectError: true,
			errorMsg:    "invalid credentials",
		},
		{
			name:        "invalid password",
			email:       testEmail,
			password:    "wrongpassword",
			expectError: true,
			errorMsg:    "invalid credentials",
		},
		{
			name:        "empty email",
			email:       "",
			password:    testPassword,
			expectError: true,
			errorMsg:    "email and password are required",
		},
		{
			name:        "empty password",
			email:       testEmail,
			password:    "",
			expectError: true,
			errorMsg:    "email and password are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := testApp.AuthService.Login(tt.email, tt.password)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
			}
		})
	}
}

func TestAuthService_GetCurrentUser(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.TeardownTestApp()

	// Create a test user
	testUser, err := testApp.CreateTestUser("current@example.com", "password123")
	require.NoError(t, err)

	tests := []struct {
		name        string
		setupCookie func(*fiber.Ctx)
		expectError bool
	}{
		{
			name: "valid user cookie",
			setupCookie: func(c *fiber.Ctx) {
				testApp.AuthService.SetUserSession(c, testUser)
			},
			expectError: false,
		},
		{
			name: "no cookie",
			setupCookie: func(c *fiber.Ctx) {
				// Don't set any cookie
			},
			expectError: true,
		},
		{
			name: "invalid user ID",
			setupCookie: func(c *fiber.Ctx) {
				c.Cookie(&fiber.Cookie{
					Name:  "user_id",
					Value: "99999", // Non-existent user
				})
			},
			expectError: true,
		},
		{
			name: "malformed cookie",
			setupCookie: func(c *fiber.Ctx) {
				c.Cookie(&fiber.Cookie{
					Name:  "user_id",
					Value: "invalid",
				})
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use the existing test app instead of creating a new context
			req, err := http.NewRequest("GET", "/test", nil)
			require.NoError(t, err)

			// Create a test handler that sets up the cookie and calls GetCurrentUser
			testApp.App.Get("/test", func(c *fiber.Ctx) error {
				// Setup cookie as specified
				tt.setupCookie(c)

				user, err := testApp.AuthService.GetCurrentUser(c)

				if tt.expectError {
					assert.Error(t, err)
					assert.Nil(t, user)
				} else {
					require.NoError(t, err)
					require.NotNil(t, user)
					assert.Equal(t, testUser.ID, user.ID)
					assert.Equal(t, testUser.Email, user.Email)
				}

				return c.SendString("OK")
			})

			resp, err := testApp.App.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()
		})
	}
}
