package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Env           string
	ServerAddress string
	StorageConn   string
	MusicServer   string
	Timeout       string
}

const (
	local = ".env.local"
	dev   = ".env.dev"
	prod  = ".env.prod"
)

func MustLoad() *Config {
	err := godotenv.Load(local)
	if err != nil {
		panic(err)
	}

	return &Config{
		Env:           os.Getenv("ENV"),
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		StorageConn:   os.Getenv("POSTGRES_CONN"),
		MusicServer:   os.Getenv("MUSIC_SERVER"),
		Timeout:       os.Getenv("TIMEOUT"),
	}
}
