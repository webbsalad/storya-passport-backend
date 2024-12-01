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

func (r *Repository) Register(ctx context.Context, name, passwordHash string) (model.Session, error) {
	var strUserID string

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.
		Insert("users").
		Columns("name", "password_hash", "created_at", "updated_at").
		Values(name, passwordHash, time.Now(), time.Now()).
		Suffix("ON CONFLICT (name) DO NOTHING RETURNING id")

	q, args, err := query.ToSql()
	if err != nil {
		return model.Session{}, fmt.Errorf("building query: %w", err)
	}

	if err = r.db.QueryRowContext(ctx, q, args...).Scan(&strUserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Session{}, model.ErrUserAlreadyExist
		}

		return model.Session{}, fmt.Errorf("execute query: %w", err)
	}

	userID, err := model.UserIDFromString(strUserID)
	if err != nil {
		return model.Session{}, fmt.Errorf("convert str to user id: %w", err)
	}

	deviceID, err := r.createSession(userID)
	if err != nil {
		return model.Session{}, fmt.Errorf("get device id: %w", err)
	}

	return model.Session{
		UserID:   userID,
		DeviceID: deviceID,
		Version:  1,
	}, nil

}

func (r *Repository) GetPasshash(ctx context.Context, name string) (string, error) {
	var passwordHash string

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.
		Select("password_hash").
		From("users").
		Where(sq.Eq{"name": name})

	q, args, err := query.ToSql()
	if err != nil {
		return "", fmt.Errorf("build query: %w", err)
	}

	if err = r.db.QueryRowContext(ctx, q, args...).Scan(&passwordHash); err != nil {
		return "", fmt.Errorf("execute query: %w", err)
	}

	return passwordHash, nil

}

func (r *Repository) GetSessionInfo(ctx context.Context, name string) (model.Session, error) {
	var strUserID string

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.
		Select("id").
		From("users").
		Where(sq.Eq{"name": name})

	q, args, err := query.ToSql()
	if err != nil {
		return model.Session{}, fmt.Errorf("build query: %w", err)
	}

	if err = r.db.QueryRowContext(ctx, q, args...).Scan(&strUserID); err != nil {
		return model.Session{}, fmt.Errorf("execute query: %w", err)
	}

	userID, err := model.UserIDFromString(strUserID)
	if err != nil {
		return model.Session{}, fmt.Errorf("convert string to user id: %w", err)
	}

	deviceID, err := r.createSession(userID)
	if err != nil {
		return model.Session{}, fmt.Errorf("get device id: %w", err)
	}

	return model.Session{
		UserID:   userID,
		DeviceID: deviceID,
		Version:  1,
	}, nil

}

func (r *Repository) createSession(userID model.UserID) (model.DeviceID, error) {
	var strDeviceID string

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.
		Insert("user_tokens").
		Columns("user_id", "version").
		Values(userID.String(), 1).
		Suffix("RETURNING device_id")

	q, args, err := query.ToSql()
	if err != nil {
		return model.DeviceID{}, fmt.Errorf("build query: %w", err)
	}

	if err = r.db.QueryRow(q, args...).Scan(&strDeviceID); err != nil {
		return model.DeviceID{}, fmt.Errorf("execute query: %w", err)
	}

	deviceID, err := model.DeviceIDFromString(strDeviceID)
	if err != nil {
		return model.DeviceID{}, fmt.Errorf("convert str to id: %w", err)
	}

	return deviceID, nil
}
