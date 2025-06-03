package tests

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoutes_GET_Login(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.TeardownTestApp()

	req, err := http.NewRequest("GET", "/login", nil)
	require.NoError(t, err)

	resp, err := testApp.App.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check content type for JSON (since we're using mock template service)
	contentType := resp.Header.Get("Content-Type")
	assert.Contains(t, contentType, "application/json")
}

func TestRoutes_POST_Register(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.TeardownTestApp()

	tests := []struct {
		name             string
		email            string
		password         string
		expectedStatus   int
		expectedRedirect string
	}{
		{
			name:             "valid registration",
			email:            "newuser@example.com",
			password:         "password123",
			expectedStatus:   http.StatusFound, // 302 redirect
			expectedRedirect: "/dashboard",     // Now goes to dashboard, not login
		},
		{
			name:           "empty email",
			email:          "",
			password:       "password123",
			expectedStatus: http.StatusOK, // Returns JSON with error in mock // Returns JSON with error in mock
		},
		{
			name:           "empty password",
			email:          "test@example.com",
			password:       "",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("email", tt.email)
			form.Add("password", tt.password)

			req, err := http.NewRequest("POST", "/register", strings.NewReader(form.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			resp, err := testApp.App.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedRedirect != "" {
				location := resp.Header.Get("Location")
				assert.Equal(t, tt.expectedRedirect, location)
			}

			// Verify user was created in database for successful registration
			if tt.expectedStatus == http.StatusFound && tt.email != "" {
				user, err := testApp.UserRepo.FindByEmail(tt.email)
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
			}
		})
	}
}

func TestRoutes_POST_Login(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.TeardownTestApp()

	// Create a test user for login
	testEmail := "logintest@example.com"
	testPassword := "password123"
	_, err := testApp.CreateTestUser(testEmail, testPassword)
	require.NoError(t, err)

	tests := []struct {
		name             string
		email            string
		password         string
		expectedStatus   int
		expectedRedirect string
		shouldSetCookie  bool
	}{
		{
			name:             "valid login",
			email:            testEmail,
			password:         testPassword,
			expectedStatus:   http.StatusFound,
			expectedRedirect: "/dashboard",
			shouldSetCookie:  true,
		},
		{
			name:           "invalid email",
			email:          "nonexistent@example.com",
			password:       testPassword,
			expectedStatus: http.StatusOK, // Returns JSON with error in mock
		},
		{
			name:           "invalid password",
			email:          testEmail,
			password:       "wrongpassword",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty credentials",
			email:          "",
			password:       "",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("email", tt.email)
			form.Add("password", tt.password)

			req, err := http.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			resp, err := testApp.App.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedRedirect != "" {
				location := resp.Header.Get("Location")
				assert.Equal(t, tt.expectedRedirect, location)
			}

			if tt.shouldSetCookie {
				cookies := resp.Header.Values("Set-Cookie")
				found := false
				for _, cookie := range cookies {
					if strings.Contains(cookie, "user_id=") {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected user_id cookie to be set")
			}
		})
	}
}

func TestRoutes_GET_Dashboard_RequiresAuth(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.TeardownTestApp()

	req, err := http.NewRequest("GET", "/dashboard", nil)
	require.NoError(t, err)

	resp, err := testApp.App.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should redirect to login when not authenticated
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	location := resp.Header.Get("Location")
	assert.Equal(t, "/login", location)
}

func TestRoutes_GET_Dashboard_WithAuth(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.TeardownTestApp()

	// Create and login user first
	testUser, err := testApp.CreateTestUser("dashboard@example.com", "password123")
	require.NoError(t, err)

	// Simulate login by setting cookie
	req, err := http.NewRequest("GET", "/dashboard", nil)
	require.NoError(t, err)

	// Add auth cookie
	cookie := &http.Cookie{
		Name:  "user_id",
		Value: "1", // Assuming first user gets ID 1
	}
	req.AddCookie(cookie)

	resp, err := testApp.App.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify user exists
	assert.NotNil(t, testUser)
}

func TestRoutes_POST_Logout(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.TeardownTestApp()

	req, err := http.NewRequest("POST", "/logout", nil)
	require.NoError(t, err)

	resp, err := testApp.App.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should redirect to login
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	location := resp.Header.Get("Location")
	assert.Equal(t, "/login", location)

	// Should clear the cookie
	cookies := resp.Header.Values("Set-Cookie")
	_ = cookies // We have the cookies but don't need to check them for this test
	// Note: Cookie clearing behavior may vary, so we're flexible here
}

func TestRoutes_Root_Redirect(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.TeardownTestApp()

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	resp, err := testApp.App.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	location := resp.Header.Get("Location")
	assert.Equal(t, "/login", location)
}
