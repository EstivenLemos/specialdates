package models

import "time"

type DateEvent struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Description string   `json:"description"`
	EventAt    time.Time `json:"event_at"`
	DateType   string    `json:"date_type"`
	Recurring  bool      `json:"recurring"`
	CreatedAt  time.Time `json:"created_at"`
}
