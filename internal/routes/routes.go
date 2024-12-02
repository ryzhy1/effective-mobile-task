package routes

import (
	"eff/internal/handlers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoutes(r *gin.Engine, songHandler *handlers.SongHandler) {
	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		songs := api.Group("/songs")
		{
			songs.POST("/", songHandler.CreateSong)
			songs.GET("/", songHandler.GetAllSongs)
			songs.GET("/:id", songHandler.GetSongByID)
			songs.PATCH("/:id", songHandler.UpdateSong)
			songs.DELETE("/:id", songHandler.DeleteSong)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
