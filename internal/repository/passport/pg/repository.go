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

func (r *Repository) GetUser(ctx context.Context, userID model.UserID) (model.User, error) {
	var user User

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.
		Select("*").
		From("users").
		Where(sq.Eq{"id": userID.String()})

	q, args, err := query.ToSql()
	if err != nil {
		return model.User{}, fmt.Errorf("build query: %w", err)
	}

	if err = r.db.GetContext(ctx, &user, q, args...); err != nil {
		return model.User{}, fmt.Errorf("execute query: %w", err)
	}

	return model.User{
		Name: user.Name,

		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *Repository) UpdateUser(ctx context.Context, userID model.UserID, name, passwordHash string) (model.User, error) {
	var user User

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	updateQuery := psql.
		Update("users").
		Set("name", name).
		Set("password_hash", passwordHash).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": userID.String()}).
		Suffix("RETURNING id, name, password_hash, created_at, updated_at")

	q, args, err := updateQuery.ToSql()
	if err != nil {
		return model.User{}, fmt.Errorf("build update query: %w", err)
	}

	if err = r.db.QueryRowxContext(ctx, q, args...).StructScan(&user); err != nil {
		return model.User{}, fmt.Errorf("execute update query: %w", err)
	}

	deleteQuery := psql.
		Delete("user_tokens").
		Where(sq.Eq{"user_id": userID.String()})

	q, args, err = deleteQuery.ToSql()
	if err != nil {
		return model.User{}, fmt.Errorf("build delete query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return model.User{}, fmt.Errorf("execute delete query: %w", err)
	}

	return model.User{
		Name: user.Name,

		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *Repository) GetPasswordHash(ctx context.Context, name string) (string, error) {
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

func (r *Repository) GetTokenVersion(ctx context.Context, userID model.UserID, deviceID model.DeviceID) (int, error) {
	var version int

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.
		Select("version").
		From("user_tokens").
		Where(sq.And{
			sq.Eq{"user_id": userID.String()},
			sq.Eq{"device_id": deviceID.String()},
		})

	q, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query: %w", err)
	}

	if err = r.db.GetContext(ctx, &version, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, model.ErrUserNotFound
		}
		return 0, fmt.Errorf("get user: %w", err)
	}

	return version, nil

}

func (r *Repository) UpdateTokenVersion(ctx context.Context, session model.Session) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.
		Update("user_tokens").
		Set("version", session.Version).
		Where(sq.And{
			sq.Eq{"user_id": session.UserID.String()},
			sq.Eq{"device_id": session.DeviceID.String()},
		})

	q, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	res, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("update session")
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affect: %w", err)
	}

	if rowAffected == 0 {
		return model.ErrDeviceNotFound
	}

	return nil

}

func (r *Repository) LogOut(ctx context.Context, userID model.UserID, deviceID model.DeviceID) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.
		Delete("user_tokens").
		Where(sq.And{
			sq.Eq{"user_id": userID.String()},
			sq.Eq{"device_id": deviceID.String()},
		})

	q, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	res, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affect: %w", err)
	}

	if rowAffected == 0 {
		return model.ErrDeviceNotFound
	}

	return nil
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
