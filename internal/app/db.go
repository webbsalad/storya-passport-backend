package app

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/webbsalad/storya-passport-backend/internal/config"
)

func initDB(cfg config.Config) *sqlx.DB {
	db, err := sqlx.Connect("postgres", cfg.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection verification failed: %v", err)
	}

	return db
}
