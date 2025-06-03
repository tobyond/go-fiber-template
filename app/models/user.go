package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"` // "-" excludes from JSON
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(email, password string) (*User, error) {
	// Validate input
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uint) (*User, error) {
	var user User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
