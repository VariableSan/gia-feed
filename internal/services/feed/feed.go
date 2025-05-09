package feed

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/VariableSan/gia-feed/internal/domain/models"
)

type Feed struct {
	log          *slog.Logger
	feedProvider FeedProvider
}

type FeedProvider interface {
	GetFeed(ctx context.Context, id string) (*models.Feed, error)
	CreateFeed(
		ctx context.Context,
		title string,
		content string,
		authorID string,
	) (string, error)
}

func New(
	log *slog.Logger,
	provider FeedProvider,
) *Feed {
	return &Feed{
		log:          log,
		feedProvider: provider,
	}
}

func (feed *Feed) GetFeed(
	ctx context.Context,
	id string,
) (*models.Feed, error) {
	const operation = "feed.GetFeed"

	log := feed.log.With(
		slog.String("operation", operation),
	)

	result, err := feed.feedProvider.GetFeed(ctx, id)
	if err != nil {
		log.Error("failed to get feed")
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("feed sent")

	return result, nil
}

func (feed *Feed) CreateFeed(
	ctx context.Context,
	title string,
	content string,
	authorID string,
) (string, error) {
	const operation = "feed.CreateFeed"

	log := feed.log.With(
		slog.String("operation", operation),
	)

	id, err := feed.feedProvider.CreateFeed(ctx, title, content, authorID)
	if err != nil {
		log.Error("failed to save feed")
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("feed saved")

	return id, nil
}
