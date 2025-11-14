package domain

import "time"

type User struct {
	ID        string
	Username  string
	TeamID    string
	IsActive  string
	CreatedAt time.Time
}
