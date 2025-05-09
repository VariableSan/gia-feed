package app

import (
	"log/slog"

	grpcapp "github.com/VariableSan/gia-feed/internal/app/grpc"
	"github.com/VariableSan/gia-feed/internal/services/feed"
	"github.com/VariableSan/gia-feed/internal/storage/sqlite"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcHost string,
	grpcPort int,
	storagePath string,
	jwtSecret string,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	feedService := feed.New(log, storage)

	grpcApp, err := grpcapp.New(log, feedService, grpcHost, grpcPort, jwtSecret)

	return &App{
		GRPCSrv: grpcApp,
	}
}
