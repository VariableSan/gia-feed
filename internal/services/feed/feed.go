package feed

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/VariableSan/gia-feed/internal/domain/models"
)

type Feed struct {
	log          *slog.Logger
	feedProvider FeedProvider
}

type FeedProvider interface {
	Feed(ctx context.Context, id string) (models.Feed, error)
	CreateFeed(
		ctx context.Context,
		title string,
		content string,
		userID string,
	) (string, error)
}

func New(
	log *slog.Logger,
	provider FeedProvider,
	tokenTTL time.Duration,
) *Feed {
	return &Feed{
		log:          log,
		feedProvider: provider,
	}
}

func (feed *Feed) CreateFeed(
	ctx context.Context,
	title string,
	content string,
	userID string,
) (string, error) {
	const operation = "feed.CreateFeed"

	log := feed.log.With(
		slog.String("operation", operation),
	)

	id, err := feed.feedProvider.CreateFeed(ctx, title, content, userID)
	if err != nil {
		log.Error("failed to save feed")
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("feed saved")

	return id, nil
}
