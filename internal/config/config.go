package config

import (
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

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
