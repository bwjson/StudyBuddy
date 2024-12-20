package main

import (
	"context"
	"fmt"
	"github.com/bwjson/StudyBuddy/configs"
	"github.com/bwjson/StudyBuddy/internal/delivery"
	"github.com/bwjson/StudyBuddy/pkg/smtp"
	"github.com/bwjson/StudyBuddy/server"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
	"time"
	"path/filepath"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialazing the logger
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	w, err := os.Getwd()

	logFilePath := filepath.Join(w, "logs", "server.txt")

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Initialiazing log files failed: %v", err)
	}

	log.SetOutput(file)
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	log.WithFields(logrus.Fields{
		"action": "start",
		"status": "success",
	}).Info("Application started successfully")

	log.Info("Starting the app...")

	cfg, err := configs.ParseConfig()
	if err != nil {
		return
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=%s",
		cfg.PostgresDB.Host,
		cfg.PostgresDB.User,
		cfg.PostgresDB.Password,
		cfg.PostgresDB.Port,
		cfg.PostgresDB.DBName,
		cfg.PostgresDB.SSLMode)

	log.Info("Starting db...")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the GORM: %v", err)
	}

	smtp := smtp.NewSMTPServer(
		cfg.SMTPServer.Host,
		cfg.SMTPServer.Port,
		cfg.SMTPServer.From,
		cfg.SMTPServer.Password)

	handler := delivery.NewHandler(db, log, smtp)

	srv := server.NewServer(cfg, handler.InitRoutes())
	go func() {
		log.Info("Starting server on port 8080...")
		srv.Run()
	}()

	log.Info("Shutting down the server...")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := srv.Stop(ctx); err != nil {
		log.Fatalf("Error while shutting down server: %s", err.Error())
	}

	log.Info("Server stopped successfully")
}
