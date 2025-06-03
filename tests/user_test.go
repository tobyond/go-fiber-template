package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
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
			name:        "valid user creation",
			email:       "test@example.com",
			password:    "password123",
			expectError: false,
		},
		{
			name:        "duplicate email",
			email:       "test@example.com", // Same as above
			password:    "password456",
			expectError: true,
		},
		{
			name:        "empty email",
			email:       "",
			password:    "password123",
			expectError: true,
		},
		{
			name:        "empty password",
			email:       "test2@example.com",
			password:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := testApp.UserRepo.Create(tt.email, tt.password)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
				assert.NotEmpty(t, user.Password)
				assert.NotEqual(t, tt.password, user.Password) // Should be hashed
				assert.Greater(t, user.ID, uint(0))
			}
		})
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.TeardownTestApp()

	// Create a test user
	createdUser, err := testApp.CreateTestUser("findme@example.com", "password123")
	require.NoError(t, err)

	tests := []struct {
		name        string
		email       string
		expectError bool
	}{
		{
			name:        "existing user",
			email:       "findme@example.com",
			expectError: false,
		},
		{
			name:        "non-existent user",
			email:       "notfound@example.com",
			expectError: true,
		},
		{
			name:        "empty email",
			email:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := testApp.UserRepo.FindByEmail(tt.email)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, createdUser.ID, user.ID)
				assert.Equal(t, tt.email, user.Email)
			}
		})
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.TeardownTestApp()

	// Create a test user
	createdUser, err := testApp.CreateTestUser("findbyid@example.com", "password123")
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      uint
		expectError bool
	}{
		{
			name:        "existing user",
			userID:      createdUser.ID,
			expectError: false,
		},
		{
			name:        "non-existent user",
			userID:      99999,
			expectError: true,
		},
		{
			name:        "zero ID",
			userID:      0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := testApp.UserRepo.FindByID(tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, createdUser.ID, user.ID)
				assert.Equal(t, createdUser.Email, user.Email)
			}
		})
	}
}

func TestUser_CheckPassword(t *testing.T) {
	testApp := SetupTestApp(t)
	defer testApp.TeardownTestApp()

	user, err := testApp.CreateTestUser("password@example.com", "correctpassword")
	require.NoError(t, err)

	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "correct password",
			password: "correctpassword",
			expected: true,
		},
		{
			name:     "incorrect password",
			password: "wrongpassword",
			expected: false,
		},
		{
			name:     "empty password",
			password: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := user.CheckPassword(tt.password)
			assert.Equal(t, tt.expected, result)
		})
	}
}
