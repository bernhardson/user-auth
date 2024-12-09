package models

import "time"

type User struct {
	ID       int64     `json:"-"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Created  time.Time `json:"created"`
}
