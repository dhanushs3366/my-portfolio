package handler

import (
	"dhanushs3366/my-portfolio/services"
	"dhanushs3366/my-portfolio/services/db"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) createBlog(c echo.Context) error {
	cookie, err := c.Cookie("auth_token")
	content := c.FormValue("content")
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

func (h *Handler) editBlog(c echo.Context) error {
	blogID := c.QueryParam("ID")
	content := c.FormValue("content")

	err := h.blogStore.EditBlog(blogID, content)
	if err != nil {
		if errors.Is(err, db.ErrNoEntityFound) {
			return c.JSON(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, "Blog edited")
}

func (h *Handler) deleteBlog(c echo.Context) error {
	blogID := c.QueryParam("ID")

	err := h.blogStore.DeleteBlog(blogID)

	if err != nil {
		if errors.Is(err, db.ErrNoEntityFound) {
			return c.JSON(http.StatusNoContent, err)
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, "Blog Deleted")
}

func (h *Handler) getBlogs(c echo.Context) error {
	blogs, err := h.blogStore.GetBlogs()
	if err != nil {
		if errors.Is(err, db.ErrNoEntityFound) {
			return c.JSON(http.StatusNoContent, err)
		}
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, blogs)
}

func (h *Handler) getBlog(c echo.Context) error {
	ID := c.Param("ID")
	blog, err := h.blogStore.GetBlogByID(ID)
	if err != nil {
		if errors.Is(err, db.ErrNoEntityFound) {
			return c.JSON(http.StatusNoContent, err)
		}
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, blog)
}
