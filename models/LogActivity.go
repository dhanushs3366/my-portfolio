package models

import (
	_ "github.com/lib/pq"
)

// uni placements are happening i need to learn sql so im raw dogging this without any orms

type LoggedActivity struct {
	Key          int `json:"all_keys"`
	MiddleClicks int `json:"middle_clicks"`
	RightClicks  int `json:"right_clicks"`
	LeftClicks   int `json:"left_clicks"`
	ExtraClicks  int `json:"extra_clicks"`
}
