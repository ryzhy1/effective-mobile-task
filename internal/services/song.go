package services

import (
	"context"
	"eff/internal/domain/dto"
	"eff/internal/domain/models"
	"eff/internal/repository"
	"eff/pkg/musicServer"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
)

type SongService struct {
	log       *slog.Logger
	repo      SongRepository
	musicInfo *musicServer.MusicInfoClient
}

type SongRepository interface {
	Create(ctx context.Context, song *models.Song) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Song, error)
	GetAll(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.Song, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewSongService(log *slog.Logger, repo SongRepository, musicInfo *musicServer.MusicInfoClient) *SongService {
	return &SongService{
		log:       log,
		repo:      repo,
		musicInfo: musicInfo,
	}
}

func (s *SongService) CreateSong(ctx context.Context, dto *dto.SongCreateDTO) (*models.Song, error) {
	const op = "song.CreateSong"

	log := s.log.With(
		slog.String("op", op),
		slog.String("group", dto.Group),
		slog.String("title", dto.Song),
	)

	log.Info("fetching song details from external API")

	songDetail, err := s.musicInfo.GetSongDetails(dto.Group, dto.Song)
	if err != nil {
		log.Error("failed to fetch song details", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("parsing to model")

	song := &models.Song{
		Group:       dto.Group,
		Title:       dto.Song,
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	log.Info("creating database row")

	if err = s.repo.Create(ctx, song); err != nil {
		log.Error("failed to create song", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("song created successfully")

	return song, nil
}

func (s *SongService) GetSongByID(ctx context.Context, id uuid.UUID) (*models.Song, error) {
	const op = "song.GetSongByID"

	log := s.log.With(
		slog.String("op", op),
		slog.String("id", id.String()),
	)

	log.Info("getting data from database")

	song, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrSongNotFound) {
			s.log.Warn("song not found", err)

			return nil, fmt.Errorf("%s: %w", op, err)
		}

		s.log.Error("failed to get song", err)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("the song is found")

	return song, nil
}

func (s *SongService) GetAllSongs(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.Song, error) {
	const op = "song.GetAllSongs"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("getting data from database")

	songs, err := s.repo.GetAll(ctx, filters, limit, offset)
	if err != nil {
		s.log.Error("failed to get song", err)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("songs found")

	return songs, nil
}

func (s *SongService) UpdateSong(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	const op = "song.UpdateSong"

	log := s.log.With(
		slog.String("op", op),
		slog.String("id", id.String()),
	)

	log.Info("updating song")

	err := s.repo.Update(ctx, id, updates)
	if err != nil {
		s.log.Error("failed to update song", err)

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("song updated successfully")

	return nil
}

func (s *SongService) DeleteSong(ctx context.Context, id uuid.UUID) error {
	const op = "song.DeleteSong"

	log := s.log.With(
		slog.String("op", op),
		slog.String("id", id.String()),
	)

	log.Info("deleting song")

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.log.Error("failed to delete song", err)

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("song deleted successfully")

	return nil
}
