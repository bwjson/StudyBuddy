package main

import (
	app "github.com/bwjson/Paypal_Microservice/internal/app"
	"github.com/bwjson/Paypal_Microservice/internal/config"
	"log/slog"
	"os"
)

func main() {
	cfg := config.ParseConfig()

	log := setupLogger(cfg.Env)

	log.Info("Start app")

	app := app.NewApp(log, cfg.GRPC.Port, cfg.StoragePath)

	app.GRPCSrv.Run()

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
