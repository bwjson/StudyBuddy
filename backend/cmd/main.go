package main

import (
	"github.com/bwjson/StudyBuddy/internal/delivery"
	"github.com/sirupsen/logrus"
)

//import (
//	"github.com/bwjson/StudyBuddy/internal/delivery"
//	"log"
//)

const addr = "10.73.62.156:8080"

// @title           StudyBuddy API
// @version         1.0
// @description     This is a sample server api.

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @host      localhost:8080
//func main() {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	client, err := paypalclient.NewClient()
//	if err != nil {
//		log.Fatalf("failed to create gRPC client: %v", err)
//	}
//
//	// Initializing the logger
//	log := logrus.New()
//	log.SetLevel(logrus.InfoLevel)
//	w, err := os.Getwd()
//	if err != nil {
//		log.Fatalf("Getting working directory failed: %v", err)
//	}
//
//	logDir := filepath.Join(w, "logs")
//	if err := os.MkdirAll(logDir, 0755); err != nil {
//		log.Fatalf("Creating logs directory failed: %v", err)
//	}
//
//	logFilePath := filepath.Join(logDir, "server.txt")
//
//	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
//	if err != nil {
//		log.Fatalf("Initialiazing log files failed: %v", err)
//	}
//
//	log.SetOutput(file)
//	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
//	log.WithFields(logrus.Fields{
//		"action": "start",
//		"status": "success",
//	}).Info("Application started successfully")
//
//	log.Info("Starting the app...")
//
//	cfg, err := configs.ParseConfig()
//	if err != nil {
//		return
//	}
//
//	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=%s",
//		cfg.PostgresDB.Host,
//		cfg.PostgresDB.User,
//		cfg.PostgresDB.Password,
//		cfg.PostgresDB.Port,
//		cfg.PostgresDB.DBName,
//		cfg.PostgresDB.SSLMode,
//	)
//
//	log.Info("Starting db...")
//
//	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//	if err != nil {
//		log.Fatalf("failed to connect to the GORM: %v", err)
//	}
//
//	smtpServ := smtp.NewSMTPServer(
//		cfg.SMTPServer.Host,
//		cfg.SMTPServer.Port,
//		cfg.SMTPServer.From,
//		cfg.SMTPServer.Password,
//	)
//
//	handler := delivery.NewHandler(db, log, smtpServ, client)
//
//	srv := server.NewServer(cfg, handler.InitRoutes())
//
//	go func() {
//		log.Info("Starting server on port 8080...")
//		srv.Run()
//	}()
//
//	log.Info("Shutting down the server...")
//
//	// Graceful shutdown
//	quit := make(chan os.Signal, 1)
//	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
//	<-quit
//
//	if err := srv.Stop(ctx); err != nil {
//		log.Fatalf("Error while shutting down server: %s", err.Error())
//	}
//
//	log.Info("Server stopped successfully")
//}

// FOR WEBSOCKETS CHAT JUST go run backend/cmd/main.go
func main() {
	wsSrv := delivery.NewWsServer(addr)
	logrus.Info("Started ws server")
	if err := wsSrv.Start(); err != nil {
		logrus.Errorf("Error with ws server: %v", err)
	}
}
