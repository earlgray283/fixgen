package models

import "time"

type Todo struct {
	ID          int64
	Title       string
	Description string
	Tags        []string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DoneAt      *time.Time
}
