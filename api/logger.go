package api

import (
	"encoding/json"
	"log"

	"github.com/labstack/echo/v4"
)

type LoggedActivity struct {
	Key          int `json:"all_keys"`
	MiddleClicks int `json:"middle_clicks"`
	RightClicks  int `json:"right_clicks"`
	LeftClicks   int `json:"left_clicks"`
	ExtraClicks  int `json:"extra_clicks"`
}

func postLogDetails(c echo.Context) error {
	var loggedActivity LoggedActivity
	err := json.NewDecoder(c.Request().Body).Decode(&loggedActivity)
	log.Print(loggedActivity)
	return err
}
