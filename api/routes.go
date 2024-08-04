package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	router *echo.Echo
	// have jwt and config future
}

type Logger struct {
	Key          int `json:"all_keys"`
	MiddleClicks int `json:"middle_clicks"`
	RightClicks  int `json:"right_clicks"`
	LeftClicks   int `json:"left_clicks"`
	ExtraClicks  int `json:"extra_clicks"`
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

	h.router.POST("/log-details", postLogDetails)

	return &h
}

func (h *Handler) Run(port uint) {
	h.router.Start(fmt.Sprintf(":%d", port))
}
