package app

import (
	"eff/internal/app/http-server"
	"eff/internal/handlers"
	"eff/internal/repository/postgres"
	"eff/internal/routes"
	"eff/internal/services"
	"eff/pkg/musicServer"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

type App struct {
	HTTPServer *http_server.Server
}

func New(log *slog.Logger, serverPort, storagePath, musicServerConn string) *App {
	storage, err := postgres.NewPostgres(storagePath)
	if err != nil {
		panic(err)
	}

	musicClient := musicServer.NewMusicInfoClient(musicServerConn)

	songService := services.NewSongService(log, storage, musicClient)
	songHandler := handlers.NewSongHandler(log, songService)

	r := gin.Default()
	err = r.SetTrustedProxies(nil)
	if err != nil {
		return nil
	}
	routes.InitRoutes(r, songHandler)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	server := http_server.NewServer(log, serverPort, r)

	return &App{
		HTTPServer: server,
	}
}
