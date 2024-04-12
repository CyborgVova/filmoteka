package entities

import (
	"time"
)

type Film struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Release     int     `json:"release"`
	Rating      int     `json:"rating"`
	Actors      []Actor `json:"actors"`
}

type Actor struct {
	ID          int       `json:"id"`
	FullName    string    `json:"fullname"`
	Sex         string    `json:"sex"`
	DateOfBirth time.Time `json:"dateofbirth"`
	Films       []Film    `json:"films"`
}
