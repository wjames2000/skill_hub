package db

import (
	"fmt"
	"time"

	"xorm.io/xorm"
	"xorm.io/xorm/log"

	"github.com/hpds/skill-hub/pkg/config"
	_ "github.com/lib/pq"
)

func Init(cfg config.DBConfig) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("xorm new engine: %w", err)
	}

	engine.SetMaxOpenConns(cfg.MaxOpen)
	engine.SetMaxIdleConns(cfg.MaxIdle)
	engine.SetConnMaxLifetime(cfg.MaxLifetime)

	engine.ShowSQL(true)
	engine.Logger().SetLevel(log.LOG_INFO)

	if err := engine.Ping(); err != nil {
		return nil, fmt.Errorf("xorm ping: %w", err)
	}

	return engine, nil
}

func WaitForDB(cfg config.DBConfig, timeout time.Duration) (*xorm.Engine, error) {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		engine, err := Init(cfg)
		if err == nil {
			return engine, nil
		}
		lastErr = err
		time.Sleep(500 * time.Millisecond)
	}
	return nil, fmt.Errorf("db wait timeout: %w", lastErr)
}
