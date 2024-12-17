package configs

import (
	"github.com/spf13/viper"
	"log"
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

func exportConfig() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config.yml")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Unable to read config: %v", err)
		return err
	}
	return nil
}

func ParseConfig() (*Config, error) {
	if err := exportConfig(); err != nil {
		return nil, err
	}

	var c Config

	if err := viper.Unmarshal(&c); err != nil {
		log.Printf("Unable to parse config: %v", err)
		return nil, err
	}

	return &c, nil
}
