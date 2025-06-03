package services

import (
	"errors"
	"fresh/app/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	userRepo *models.UserRepository
}

func NewAuthService(userRepo *models.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Login(email, password string) (*models.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.CheckPassword(password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *AuthService) Register(email, password string) (*models.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	// Check if user already exists
	if _, err := s.userRepo.FindByEmail(email); err == nil {
		return nil, errors.New("email already exists")
	}

	return s.userRepo.Create(email, password)
}

func (s *AuthService) GetCurrentUser(c *fiber.Ctx) (*models.User, error) {
	userIDStr := c.Cookies("user_id")
	if userIDStr == "" {
		return nil, errors.New("not authenticated")
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		// Clear invalid cookie
		c.ClearCookie("user_id")
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.FindByID(uint(userID))
	if err != nil {
		// Clear cookie for non-existent user
		c.ClearCookie("user_id")
		return nil, err
	}

	return user, nil
}

func (s *AuthService) SetUserSession(c *fiber.Ctx, user *models.User) {
	c.Cookie(&fiber.Cookie{
		Name:  "user_id",
		Value: strconv.Itoa(int(user.ID)),
		// Add security options for production:
		// HTTPOnly: true,
		// Secure: true,
		// SameSite: "Strict",
	})
}

func (s *AuthService) ClearUserSession(c *fiber.Ctx) {
	c.ClearCookie("user_id")
}
