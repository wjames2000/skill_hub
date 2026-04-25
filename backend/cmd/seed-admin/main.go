package main

import (
	"flag"
	"fmt"

	"github.com/hpds/skill-hub/pkg/config"
	"github.com/hpds/skill-hub/pkg/db"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	username := flag.String("username", "admin", "admin username")
	email := flag.String("email", "admin@skillhub.dev", "admin email")
	password := flag.String("password", "", "admin password")
	flag.Parse()

	if *password == "" {
		fmt.Println("error: --password is required")
		flag.Usage()
		return
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic("load config: " + err.Error())
	}

	dbEngine, err := db.Init(cfg.DB)
	if err != nil {
		panic("db init: " + err.Error())
	}
	defer dbEngine.Close()

	if err := db.Migrate(dbEngine, "migrations"); err != nil {
		panic("db migrate: " + err.Error())
	}
	fmt.Println("migrations applied")

	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		panic("bcrypt hash: " + err.Error())
	}

	rows, err := dbEngine.Query("SELECT 1 FROM users WHERE username = $1 OR email = $2 LIMIT 1", *username, *email)
	if err != nil {
		panic("check user: " + err.Error())
	}
	if len(rows) > 0 {
		fmt.Println("user already exists, skipping creation")
		return
	}

	_, err = dbEngine.Exec(
		`INSERT INTO users (username, email, password_hash, role, status) VALUES ($1, $2, $3, 'admin', 1)`,
		*username, *email, string(hash),
	)
	if err != nil {
		panic("insert admin: " + err.Error())
	}

	fmt.Printf("admin user created: username=%s email=%s role=admin\n", *username, *email)
}
