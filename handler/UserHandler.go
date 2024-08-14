package handler

import (
	"dhanushs3366/my-portfolio/services"
	"dhanushs3366/my-portfolio/services/db"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Login(c echo.Context) error {
	req := c.Request()

	// dont add bcrypt right away i hardcoded my initial creds without hashing it
	// use the update route to update it into hashed

	username := req.FormValue("username")
	password := req.FormValue("password")

	user, err := h.userStore.GetUser(username)

	if err != nil {
		if errors.Is(err, db.ErrNoEntityFound) {
			return c.JSON(http.StatusUnauthorized, "user not found")
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
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
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(services.EXPIRY_TIME * time.Hour),
	})

	return c.JSON(http.StatusOK, "login successful")
}

func (h *Handler) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, "logged out")
}

func (h *Handler) UpdatePassword(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.FormValue("password")

	hashedPassword, err := services.HashPassword(password)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = h.userStore.UpdatePassword(username, hashedPassword)

	if err != nil {
		if errors.Is(err, db.ErrNoEntityFound) {
			return c.JSON(http.StatusUnauthorized, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "updated password sucessfully")
}
