package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/VariableSan/gia-feed/internal/domain/models"
	"github.com/VariableSan/gia-feed/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const operation = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) SaveFeed(ctx context.Context, title string, content string) (int64, error) {
	const operation = "storage.sqlite.SaveFeed"

	stmt, err := s.db.Prepare("INSERT INTO feeds(title, content) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	res, err := stmt.ExecContext(ctx, title, content)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", operation, storage.ErrFeedExists)
		}

		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	return id, nil
}

func (s *Storage) Feed(ctx context.Context, uid int64) (models.Feed, error) {
	const operation = "storage.sqlite.Feed"

	stmt, err := s.db.Prepare("SELECT id, title, content FROM feeds WHERE id = ?")
	if err != nil {
		return models.Feed{}, fmt.Errorf("%s: %w", operation, err)
	}

	row := stmt.QueryRowContext(ctx, uid)

	var feed models.Feed
	err = row.Scan(&feed.ID, &feed.Content, &feed.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Feed{}, fmt.Errorf("%s: %w", operation, storage.ErrFeedNotFound)
		}

		return models.Feed{}, fmt.Errorf("%s: %w", operation, err)
	}

	return feed, nil
}
