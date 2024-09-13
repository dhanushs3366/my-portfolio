package handler

import (
	"dhanushs3366/my-portfolio/api"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func (h *Handler) getRepos(c echo.Context) error {
	username := os.Getenv("GITHUB_USERNAME")
	repos, err := api.FetchReposByUserName(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, repos)
}
