package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App     AppConfig     `mapstructure:"app"`
	DB      DBConfig      `mapstructure:"db"`
	Redis   RedisConfig   `mapstructure:"redis"`
	Meili   MeiliConfig   `mapstructure:"meili"`
	Minio   MinioConfig   `mapstructure:"minio"`
	Consul  ConsulConfig  `mapstructure:"consul"`
	GitHub  GitHubConfig  `mapstructure:"github"`
	Sync    SyncConfig    `mapstructure:"sync"`
	Asynq   AsynqConfig   `mapstructure:"asynq"`
	Semgrep SemgrepConfig `mapstructure:"semgrep"`
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

type GitHubConfig struct {
	Tokens       []string `mapstructure:"tokens"`
	SearchTopics []string `mapstructure:"search_topics"`
	MaxPerPage   int      `mapstructure:"max_per_page"`
	RequestDelay int      `mapstructure:"request_delay"`
}

type SyncConfig struct {
	FullSyncCron    string `mapstructure:"full_sync_cron"`
	IncrementalCron string `mapstructure:"incremental_cron"`
	IncrementalDays int    `mapstructure:"incremental_days"`
	Concurrency     int    `mapstructure:"concurrency"`
	SyncTimeout     int    `mapstructure:"sync_timeout"`
	ScanEnabled     bool   `mapstructure:"scan_enabled"`
}

type AsynqConfig struct {
	Enabled     bool `mapstructure:"enabled"`
	Concurrency int  `mapstructure:"concurrency"`
}

type SemgrepConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Binary  string `mapstructure:"binary"`
	Rules   string `mapstructure:"rules"`
	Timeout int    `mapstructure:"timeout"`
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
