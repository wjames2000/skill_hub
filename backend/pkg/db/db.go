package db

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
	ensureDB(cfg)

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

func Migrate(engine *xorm.Engine, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read migration dir: %w", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, f := range files {
		path := filepath.Join(dir, f)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", f, err)
		}

		if _, err := engine.Exec(string(content)); err != nil {
			return fmt.Errorf("exec %s: %w", f, err)
		}
	}

	return nil
}

func ensureDB(cfg config.DBConfig) {
	dbName := extractDBName(cfg.DSN)
	if dbName == "" {
		return
	}

	adminDSN := strings.Replace(cfg.DSN, fmt.Sprintf("dbname=%s", dbName), "dbname=postgres", 1)
	adminDSN = strings.Replace(adminDSN, fmt.Sprintf("database=%s", dbName), "database=postgres", 1)

	engine, err := xorm.NewEngine(cfg.Driver, adminDSN)
	if err != nil {
		return
	}
	defer engine.Close()

	exists := false
	rows, err := engine.Query(fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", dbName))
	if err == nil && len(rows) > 0 {
		exists = true
	}

	if !exists {
		_, _ = engine.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	}
}

func extractDBName(dsn string) string {
	for _, part := range strings.Split(dsn, " ") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "dbname=") {
			return strings.TrimPrefix(part, "dbname=")
		}
		if strings.HasPrefix(part, "database=") {
			return strings.TrimPrefix(part, "database=")
		}
	}
	return ""
}
