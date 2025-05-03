package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/VariableSan/gia-feed/internal/app/grpc"
	"github.com/VariableSan/gia-feed/internal/services/feed"
	"github.com/VariableSan/gia-feed/internal/storage/sqlite"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	feedService := feed.New(log, storage, tokenTTL)

	grpcApp := grpcapp.New(log, feedService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
