package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string     `yaml:"env" env-default:"local"`
	StoragePath string     `yaml:"storage_path" env-required:"./data"`
	GRPC        GRPCConfig `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func ParseConfig() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config file path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	// --config="path/to/config.yml"
	flag.StringVar(&res, "config", "", "path to config")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
