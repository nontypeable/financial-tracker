package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Database     *DatabaseConfig     `mapstructure:"database"`
		Server       *ServerConfig       `mapstructure:"server"`
		TokenManager *TokenManagerConfig `mapstructure:"token_manager"`
	}

	ServerConfig struct {
		Address         string        `mapstructure:"address"`
		ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	}

	TokenManagerConfig struct {
		AccessSecret  string        `mapstructure:"access_secret"`
		RefreshSecret string        `mapstructure:"refresh_secret"`
		AccessTTL     time.Duration `mapstructure:"access_ttl"`
		RefreshTTL    time.Duration `mapstructure:"refresh_ttl"`
	}

	DatabaseConfig struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"db_name"`
		SSLMode  string `mapstructure:"ssl_mode"`
	}
)

var (
	cfg     *Config
	once    sync.Once
	loadErr error
)

func LoadConfig(path string) (*Config, error) {
	once.Do(func() {
		v := viper.New()
		v.SetConfigFile(path)
		v.SetConfigType("yaml")

		if err := v.ReadInConfig(); err != nil {
			loadErr = fmt.Errorf("failed to read config file: %w", err)
			return
		}

		var c Config
		if err := v.Unmarshal(&c); err != nil {
			loadErr = fmt.Errorf("failed to unmarshal config: %w", err)
			return
		}

		cfg = &c
	})

	return cfg, loadErr
}
