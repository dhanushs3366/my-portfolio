package services

import (
	"dhanushs3366/my-portfolio/models"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	HASHING_ROUNDS = 14
	EXPIRY_TIME    = 25 // in hrs
)

type UserClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), HASHING_ROUNDS)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func GenerateJWTToken(user *models.User) (string, error) {
	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return "", errors.New("no env var found")
	}

	expirationTime := time.Now().Add(EXPIRY_TIME * time.Hour)

	claims := &UserClaims{
		UserID:   user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func ValidateUser()
