package pg

import (
	"github.com/jmoiron/sqlx"
	"github.com/webbsalad/storya-passport-backend/internal/repository/passport"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) (passport.Repository, error) {
	return &Repository{db: db}, nil
}
