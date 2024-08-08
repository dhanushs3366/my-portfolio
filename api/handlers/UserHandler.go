package handlers

import (
	"dhanushs3366/my-portfolio/services"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	req := c.Request()

	// dont add bcrypt right away i hardcoded my initial creds without hashing it
	// use the update route to update it into hashed

	username := req.FormValue("username")
	password := req.FormValue("password")

	user, err := services.GetUser(username)

	if err != nil {
		if errors.Is(err, services.ErrNoEntityFound) {
			return c.JSON(http.StatusUnauthorized, "user not found")
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if password != user.Password {
		return c.JSON(http.StatusUnauthorized, "incorrect username/password")
	}

	if !user.IsAdmin {
		return c.JSON(http.StatusUnauthorized, "unauthorised user")
	}

	// user is now admin and authorised generate token and set it in cookie
	token, err := services.GenerateJWTToken(user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.SetCookie(&http.Cookie{
		Name:     "Bearer",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(services.EXPIRY_TIME * time.Hour),
	})

	return c.JSON(http.StatusOK, "login successful")
}
