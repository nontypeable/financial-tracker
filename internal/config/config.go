package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Database *DatabaseConfig `yaml:"database"`
	}

	DatabaseConfig struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"db_name"`
		SSLMode  string `yaml:"ssl_mode"`
	}
)

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("$CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file is not exists")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read config")
	}

	return &cfg
}
