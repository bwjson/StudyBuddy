package grpcApp

import (
	"fmt"
	grpcSrv "github.com/bwjson/Paypal_Microservice/internal/grpc"
	"github.com/bwjson/Paypal_Microservice/storage/sqlite"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log  *slog.Logger
	gRPC *grpc.Server
	port string
}

func NewApp(log *slog.Logger, storage *sqlite.Storage, port string) *App {
	gRPCServer := grpc.NewServer()

	grpcSrv.Register(gRPCServer, storage)

	return &App{
		log:  log,
		gRPC: gRPCServer,
		port: port,
	}
}

func (app *App) Run() error {
	const op = "grpcApp.Run"

	log := app.log.With(slog.String("op", op))

	l, err := net.Listen("tcp", ":44044")
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("%s: %w", op, err.Error())
	}

	log.Info("Starting gRPC server", slog.String("address", l.Addr().String()))

	if err := app.gRPC.Serve(l); err != nil {
		log.Error(err.Error())
		return fmt.Errorf("%s: %w", op, err.Error())
	}

	return nil
}

func (app *App) Stop() {
	const op = "grpcApp.Stop"

	app.log.With(slog.String("op", op)).Info("Stopping gRPC server", slog.String("port", app.port))

	app.gRPC.GracefulStop()
}
