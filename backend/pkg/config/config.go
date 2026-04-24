package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	DB     DBConfig     `mapstructure:"db"`
	Redis  RedisConfig  `mapstructure:"redis"`
	Meili  MeiliConfig  `mapstructure:"meili"`
	Minio  MinioConfig  `mapstructure:"minio"`
	Consul ConsulConfig `mapstructure:"consul"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Env     string `mapstructure:"env"`
	Port    int    `mapstructure:"port"`
	Version string `mapstructure:"version"`
}

type DBConfig struct {
	Driver      string        `mapstructure:"driver"`
	DSN         string        `mapstructure:"dsn"`
	MaxOpen     int           `mapstructure:"max_open"`
	MaxIdle     int           `mapstructure:"max_idle"`
	MaxLifetime time.Duration `mapstructure:"max_lifetime"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type MeiliConfig struct {
	Host   string `mapstructure:"host"`
	APIKey string `mapstructure:"api_key"`
}

type MinioConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	UseSSL    bool   `mapstructure:"use_ssl"`
	Bucket    string `mapstructure:"bucket"`
}

type ConsulConfig struct {
	Addr  string `mapstructure:"addr"`
	Token string `mapstructure:"token"`
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.AutomaticEnv()
	v.SetEnvPrefix("SKILL_HUB")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
