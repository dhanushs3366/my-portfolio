package handler

import (
	"dhanushs3366/my-portfolio/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) createBlog(c echo.Context) error {
	cookie, err := c.Cookie("auth_token")
	content := c.FormValue("blog-content")
	if err != nil {
		return c.JSON(http.StatusNoContent, http.ErrNoCookie)
	}

	username, err := services.GetUsernameFromToken(cookie.Value)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	user, err := h.userStore.GetUser(username)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	err = h.blogStore.CreateBlog(user, content)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, "Blog created")
}
