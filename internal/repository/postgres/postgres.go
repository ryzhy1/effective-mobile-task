package postgres

import (
	"context"
	"eff/internal/domain/models"
	"eff/internal/repository"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewPostgres(conn string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := pgxpool.New(context.Background(), conn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (r *Storage) Create(ctx context.Context, song *models.Song) error {
	const op = "storage.Create"

	query, args, err := squirrel.Insert("songs").
		Columns("group_name", "title", "release_date", "text", "link", "created_at").
		Values(song.Group, song.Title, song.ReleaseDate, song.Text, song.Link, time.Now()).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Storage) GetByID(ctx context.Context, id uuid.UUID) (*models.Song, error) {
	const op = "storage.GetByID"

	query, args, err := squirrel.Select("group_name", "title", "release_date", "text", "link").
		From("songs").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var song models.Song
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&song.Group, &song.Title, &song.ReleaseDate, &song.Text, &song.Link,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &song, nil
}

func (r *Storage) GetAll(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.Song, error) {
	const op = "storage.GetAll"

	queryBuilder := squirrel.Select("group_name", "title", "release_date", "text", "link").
		From("songs").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	for key, value := range filters {
		queryBuilder = queryBuilder.Where(squirrel.Eq{key: value})
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(
			&song.Group, &song.Title, &song.ReleaseDate, &song.Text, &song.Link,
		); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("%s: %w", op, repository.ErrSongNotFound)
			}
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func (r *Storage) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	const op = "storage.Update"

	query, args, err := squirrel.Update("songs").
		SetMap(updates).
		SetMap(squirrel.Eq{"updated_at": time.Now()}).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Storage) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "storage.Delete"

	query, args, err := squirrel.Delete("songs").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
