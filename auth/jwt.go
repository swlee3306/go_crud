package auth

import (
	"errors"
	"fmt"
	"go-crud/config"
	"go-crud/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type AuthService struct {
	secretKey []byte
}

func NewAuthService(secretKey string) *AuthService {
	return &AuthService{
		secretKey: []byte(secretKey),
	}
}

func (as *AuthService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func (as *AuthService) CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (as *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(as.secretKey)
}

func (as *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return as.secretKey, nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("invalid token")
}

func (as *AuthService) RefreshToken(tokenString string) (string, error) {
	claims, err := as.ValidateToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("failed to validate token: %w", err)
	}
	
	// Check if token is close to expiry (within 1 hour)
	if time.Until(claims.ExpiresAt.Time) > time.Hour {
		return "", errors.New("token is not close to expiry")
	}
	
	// Generate new token with same claims but new expiry
	user := &models.User{
		ID:       claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
	}
	
	return as.GenerateToken(user)
}

func (as *AuthService) Login(email, password string) (*models.User, string, error) {
	var user models.User
	
	// Find user by email
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}
	
	// Check if user is active
	if !user.IsActive {
		return nil, "", fmt.Errorf("account is deactivated")
	}
	
	// Verify password
	if err := as.CheckPassword(user.Password, password); err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}
	
	// Generate token
	token, err := as.GenerateToken(&user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}
	
	return &user, token, nil
}

func (as *AuthService) Register(userData *models.User) (*models.User, string, error) {
	// Hash password
	hashedPassword, err := as.HashPassword(userData.Password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}
	userData.Password = hashedPassword
	
	// Set default values
	userData.IsActive = true
	
	// Create user
	if err := config.DB.Create(userData).Error; err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}
	
	// Generate token
	token, err := as.GenerateToken(userData)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}
	
	return userData, token, nil
}
