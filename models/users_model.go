package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Nickname string    `json:"nickname"`
	Img      string    `json:"img"`
	Country  string    `json:"country"`
	City     string    `json:"city"`
	Clubs    []Club    `json:"clubs"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Deleted  bool      `json:"deleted"`
}

type Club struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
