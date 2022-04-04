package model

import "time"

type User struct {
	ID        int64
	Username  string
	Password  string
	Email     string
	Phone     string
	Address   string
	Roles     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
