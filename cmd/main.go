package main

import (
	"github.com/bwjson/StudyBuddy/configs"
	"github.com/bwjson/StudyBuddy/internal/delivery"
	_ "github.com/bwjson/StudyBuddy/pkg/postgres"
	"github.com/bwjson/StudyBuddy/server"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	// TODO: Инициализация логера

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	cfg, err := configs.ParseConfig()
	if err != nil {
		return
	}

	//postgres, err := postgres.NewPostgresDB(*cfg)
	//if err != nil {
	//	return
	//}

	// TODO: Заменить хардкод на конфиг
	dsn := "host=localhost user=postgres password=5432 dbname=studybuddy port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the GORM: %v", err)
	}

	handler := delivery.NewHandler(db)

	srv := server.NewServer(cfg, handler.InitRoutes())

	srv.Run()
}
