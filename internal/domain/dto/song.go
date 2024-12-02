package dto

import (
	"github.com/google/uuid"
)

type Song struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Group       string    `json:"group"`
	Title       string    `json:"title"`
	ReleaseDate string    `json:"release_date"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}
