package handlers

import (
	"dhanushs3366/my-portfolio/api/github"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetRepos(c echo.Context) error {
	username := c.QueryParam("username")
	repos, err := github.FetchReposByUserName(username)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, repos)
}

func GetGitUser(c echo.Context) error {
	user, err := github.FetchUser()

	if err != nil {
		return c.JSON(http.StatusNotFound, "user not found")
	}

	return c.JSON(http.StatusOK, user)
}
