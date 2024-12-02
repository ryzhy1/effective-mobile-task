package repository

import (
	"errors"
)

var (
	ErrSongNotFound      = errors.New("song not found")
	ErrSongAlreadyExists = errors.New("song already exists")
)
