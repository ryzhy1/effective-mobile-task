package dto

type SongCreateDTO struct {
	Group string `json:"group" binding:"required"`
	Song  string `json:"song" binding:"required"`
}
