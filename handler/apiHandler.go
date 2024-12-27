package handler

import (
	"dhanushs3366/my-portfolio/services/hooks"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) fetchGitRepos(c echo.Context) error {
	repos, err := hooks.FetchGitRepos()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, repos)
}
