package models

import (
	"github.com/google/uuid"
	"time"
)

type Song struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Group       string    `json:"group"`
	Title       string    `json:"title"`
	ReleaseDate string    `json:"release_date"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
