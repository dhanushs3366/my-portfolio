package services

import (
	"dhanushs3366/my-portfolio/models"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
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

func ValidateJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		cookie, err := c.Cookie("auth_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				return c.JSON(http.StatusUnauthorized, "no cookie")
			}
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		tokenStr := cookie.Value

		claims := &UserClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				return c.JSON(http.StatusUnauthorized, err.Error())
			}
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		if !token.Valid {
			return c.JSON(http.StatusUnauthorized, "user unauthorized")
		}

		return next(c)

	}
}
