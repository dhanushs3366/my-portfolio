package models

import "time"

type Blog struct {
	ID      int
	Date    time.Time
	Content string
	OwnedBy int
}
