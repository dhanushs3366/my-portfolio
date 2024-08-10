package models

import "time"

type LoggedActivity struct {
	ID           int       `json:"ID"`
	Key          int       `json:"all_keys"`
	MiddleClicks int       `json:"middle_clicks"`
	RightClicks  int       `json:"right_clicks"`
	LeftClicks   int       `json:"left_clicks"`
	ExtraClicks  int       `json:"extra_clicks"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
