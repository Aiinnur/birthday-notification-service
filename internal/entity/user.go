package entity

import "time"

type User struct {
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Birthday time.Time `json:"birthday"`
}
