package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/webbsalad/storya-passport-backend/internal/model"
	"github.com/webbsalad/storya-passport-backend/internal/repository/passport"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) (passport.Repository, error) {
	return &Repository{db: db}, nil
}

func (r *Repository) Register(ctx context.Context, name, passwordHash string) (model.UserID, error) {
	var strUserID string
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.
		Insert("users").
		Columns("name", "password_hash", "created_at", "updated_at").
		Values(name, passwordHash, time.Now(), time.Now()).
		Suffix("ON CONFLICT (name) DO NOTHING RETURNING id")

	q, args, err := query.ToSql()
	if err != nil {
		return model.UserID{}, fmt.Errorf("building query: %w", err)
	}

	err = r.db.QueryRowContext(ctx, q, args...).Scan(&strUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.UserID{}, model.ErrUserAlreadyExist
		}

		return model.UserID{}, fmt.Errorf("execute query: %w", err)
	}

	userID, err := model.UserIDFromString(strUserID)
	if err != nil {
		return model.UserID{}, fmt.Errorf("convert str to user id: %w", err)
	}

	return userID, err

}
