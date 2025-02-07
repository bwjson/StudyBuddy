package app

import (
	grpcApp "github.com/bwjson/Paypal_Microservice/internal/app/grpc"
	"github.com/bwjson/Paypal_Microservice/storage/sqlite"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcApp.App
}

func NewApp(log *slog.Logger, port string, storagePath string) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	grpcApp := grpcApp.NewApp(log, storage, port)

	return &App{
		GRPCSrv: grpcApp,
	}
}
