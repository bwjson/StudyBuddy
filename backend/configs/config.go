package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server     Server
	PostgresDB PostgresDB
	SMTPServer SMTPServer
}

type Server struct {
	Port               string
	Development        bool
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	MaxHeaderMegabytes int
}

type PostgresDB struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type SMTPServer struct {
	Host     string
	Port     string
	From     string
	Password string
}

func ParseConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
		return nil, err
	}

	readTimeout, _ := time.ParseDuration(os.Getenv("SERVER_READ_TIMEOUT"))
	writeTimeout, _ := time.ParseDuration(os.Getenv("SERVER_WRITE_TIMEOUT"))
	maxHeaderMegabytes, _ := strconv.Atoi(os.Getenv("SERVER_MAX_HEADER_MEGABYTES"))

	development, _ := strconv.ParseBool(os.Getenv("SERVER_DEVELOPMENT"))

	c := &Config{
		Server: Server{
			Port:               os.Getenv("SERVER_PORT"),
			Development:        development,
			ReadTimeout:        readTimeout,
			WriteTimeout:       writeTimeout,
			MaxHeaderMegabytes: maxHeaderMegabytes,
		},
		PostgresDB: PostgresDB{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   os.Getenv("POSTGRES_DB_NAME"),
			SSLMode:  os.Getenv("POSTGRES_SSL_MODE"),
		},
		SMTPServer: SMTPServer{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     os.Getenv("SMTP_PORT"),
			From:     os.Getenv("SMTP_FROM"),
			Password: os.Getenv("SMTP_PASSWORD"),
		},
	}

	return c, nil
}

//func exportConfig() error {
//	viper.SetConfigType("yaml")
//	viper.AddConfigPath("./configs")
//	viper.SetConfigName("config.yml")
//
//	if err := viper.ReadInConfig(); err != nil {
//		log.Printf("Unable to read config: %v", err)
//		return err
//	}
//	return nil
//}
//
//func ParseConfig() (*Config, error) {
//	if err := exportConfig(); err != nil {
//		return nil, err
//	}
//
//	var c Config
//
//	if err := viper.Unmarshal(&c); err != nil {
//		log.Printf("Unable to parse config: %v", err)
//		return nil, err
//	}
//
//	return &c, nil
//}
