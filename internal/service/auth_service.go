package service

import (
	"errors"
	"os"
	"time"

	"github.com/Richa3822/recipe-api/internal/model"
	"github.com/Richa3822/recipe-api/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req *model.RegisterRequest) (*model.AuthResponse, error)
	Login(req *model.LoginRequest) (*model.AuthResponse, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo: repo}
}

// REGISTER
func (s *authService) Register(req *model.RegisterRequest) (*model.AuthResponse, error) {
	// Check if email already exists
	_, err := s.repo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already registered")
	}

	// Hash the password — never store plain text
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate JWT token
	token, err := generateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{Token: token, User: *user}, nil
}

// LOGIN
func (s *authService) Login(req *model.LoginRequest) (*model.AuthResponse, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password") // same message — don't leak which field is wrong
	}

	// Generate JWT token
	token, err := generateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{Token: token, User: *user}, nil
}

// JWT token generator — private helper function
func generateToken(user *model.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	// Claims = payload of the token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
