package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/VariableSan/gia-feed/internal/app"
	"github.com/VariableSan/gia-feed/internal/config"
	"github.com/VariableSan/gia-feed/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.NewLogger(cfg.Env)

	log.Info(
		"starting application on port: ",
		slog.Int("port", cfg.GRPC.Port),
	)

	application := app.New(
		log,
		cfg.GRPC.Host,
		cfg.GRPC.Port,
		cfg.StoragePath,
		cfg.JWTSecret,
	)

	go application.GRPCSrv.MustRun()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case stopSignal := <-interrupt:
		log.Info(
			"stopping application",
			slog.String("signal", stopSignal.String()),
		)
	}

	application.GRPCSrv.Stop()

	log.Info("application stopped")
}
