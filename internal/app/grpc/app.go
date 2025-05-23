package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	feedgrpc "github.com/VariableSan/gia-feed/internal/grpc/feed"
	"github.com/VariableSan/gia-feed/internal/grpc/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	host       string
	port       int
	jwtSecret  string
}

func New(
	log *slog.Logger,
	feedService feedgrpc.FeedService,
	host string,
	port int,
	jwtSecret string,
) (*App, error) {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthInterceptor([]byte(jwtSecret))),
	)

	feedgrpc.Register(gRPCServer, feedService)
	reflection.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		host:       host,
		port:       port,
	}, nil
}

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *App) Run() error {
	const operation = "grpcapp.Run"

	log := app.log.With(
		slog.String("operation", operation),
		slog.Int("port", app.port),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", app.host, app.port))
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("grpc server is running", slog.String("addr", listener.Addr().String()))

	if err := app.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}

func (app *App) Stop() {
	const operation = "grpcapp.Stop"

	app.log.
		With(slog.String("operation", operation)).
		Info("stopping gRPC server", slog.Int("port", app.port))

	app.gRPCServer.GracefulStop()
}
