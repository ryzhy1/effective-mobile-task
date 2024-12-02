package handlers

import (
	"eff/internal/domain/dto"
	"eff/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"strconv"
)

// ErrorResponse структура для ошибок и успешных сообщений
// @Description Структура ответа с ошибкой или успешным сообщением
type ErrorResponse struct {
	Message string `json:"message"`         // сообщение об ошибке
	Error   string `json:"error,omitempty"` // подробности ошибки
}

type SongHandler struct {
	log     *slog.Logger
	service *services.SongService
}

func NewSongHandler(log *slog.Logger, service *services.SongService) *SongHandler {
	return &SongHandler{
		log:     log,
		service: service,
	}
}

// CreateSong создаёт новую песню
// @Summary Create a song
// @Description Create a new song in the system
// @Tags songs
// @Accept json
// @Produce json
// @Param song body dto.SongCreateDTO true "Song data"
// @Success 200 {object} dto.Song
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/songs [post]
func (h *SongHandler) CreateSong(c *gin.Context) {
	var input dto.SongCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	song, err := h.service.CreateSong(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, song)
}

// GetSongByID получает песню по её ID
// @Summary Get a song by ID
// @Description Retrieve a song from the system using its ID
// @Tags songs
// @Produce json
// @Param id path string true "Song ID"
// @Success 200 {object} dto.Song
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/songs/{id} [get]
func (h *SongHandler) GetSongByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid song ID"})
		return
	}

	song, err := h.service.GetSongByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	if song == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Song not found"})
		return
	}

	c.JSON(http.StatusOK, song)
}

// GetAllSongs получает все песни с фильтрами и пагинацией
// @Summary Get all songs
// @Description Retrieve a list of songs with optional filters and pagination
// @Tags songs
// @Produce json
// @Param group query string false "Group"
// @Param title query string false "Title"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} dto.Song
// @Failure 500 {object} ErrorResponse
// @Router /api/songs [get]
func (h *SongHandler) GetAllSongs(c *gin.Context) {
	filters := map[string]interface{}{}
	if group := c.Query("group"); group != "" {
		filters["group"] = group
	}
	if title := c.Query("title"); title != "" {
		filters["title"] = title
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	songs, err := h.service.GetAllSongs(c.Request.Context(), filters, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, songs)
}

// UpdateSong обновляет информацию о песне
// @Summary Update a song
// @Description Update an existing song's details
// @Tags songs
// @Accept json
// @Produce json
// @Param id path string true "Song ID"
// @Param song body map[string]interface{} true "Updated song data"
// @Success 200 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/songs/{id} [patch]
func (h *SongHandler) UpdateSong(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid song ID"})
		return
	}

	updates := map[string]interface{}{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.UpdateSong(c.Request.Context(), id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ErrorResponse{Message: "Song updated successfully"})
}

// DeleteSong удаляет песню
// @Summary Delete a song
// @Description Delete a song from the system by ID
// @Tags songs
// @Produce json
// @Param id path string true "Song ID"
// @Success 200 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/songs/{id} [delete]
func (h *SongHandler) DeleteSong(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid song ID"})
		return
	}

	if err := h.service.DeleteSong(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ErrorResponse{Message: "Song deleted successfully"})
}
