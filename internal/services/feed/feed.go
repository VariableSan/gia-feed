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
	Feed(ctx context.Context, uid int64) (models.Feed, error)
	SaveFeed(
		ctx context.Context,
		title string,
		content string,
	) (int64, error)
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

func (feed *Feed) SaveFeed(
	ctx context.Context,
	title string,
	content string,
) (int64, error) {
	const operation = "feed.SaveFeed"

	log := feed.log.With(
		slog.String("operation", operation),
	)

	id, err := feed.feedProvider.SaveFeed(ctx, title, content)
	if err != nil {
		log.Error("failed to save feed")
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("feed saved")

	return id, nil
}
