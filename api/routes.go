package api

import (
	"dhanushs3366/my-portfolio/api/handlers"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	router *echo.Echo
	// have jwt and config future
}

type e echo.HandlerFunc

func writeJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func Init() *Handler {
	h := Handler{router: echo.New()}
	h.router.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "HIIII")
	})

	h.router.POST("/log-details", handlers.PostLogDetails)

	h.router.GET("/repos", handlers.GetRepos)
	h.router.GET("/git-user", handlers.GetGitUser)

	return &h
}

func (h *Handler) Run(port uint) {
	h.router.Start(fmt.Sprintf(":%d", port))
}
