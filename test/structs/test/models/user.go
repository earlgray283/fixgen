package models

import "time"

type User struct {
	ID        int64
	Name      string
	IconURL   string
	UserType  int64
	CreatedAt time.Time
	UpdatedAt *time.Time
}
