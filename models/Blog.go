package models

import "time"

type Blog struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	OwnedBy   int       `json:"owned-by"`
	CreatedAt time.Time `json:"created-at"`
	UpdatedAt time.Time `json:"updated-at"`
}
