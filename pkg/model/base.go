package model

import "time"

type Base struct {
	ID        uint64    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Entry struct {
	Name       string
	Properties map[string]string
}
