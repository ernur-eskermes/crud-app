package psql

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/jackc/pgconn"

	"github.com/ernur-eskermes/crud-app/internal/core"
	"github.com/ernur-eskermes/crud-app/pkg/database/postgresql"
	"github.com/jackc/pgx/v4"
)

type UsersRepo struct {
	db postgresql.Client
}

func NewUsersRepo(db postgresql.Client) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) GetByID(ctx context.Context, id uuid.UUID) (core.User, error) {
	q := "SELECT id, username, password FROM users WHERE id=$1"

	var user core.User

	if err := r.db.QueryRow(ctx, q, id).Scan(&user.ID, &user.Username, &user.Password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core.User{}, core.ErrUserNotFound
		}

		return core.User{}, err
	}

	return user, nil
}

func (r *UsersRepo) Verify(ctx context.Context, username string) error {
	q := "UPDATE users SET is_active=true WHERE username=$1"

	if res, err := r.db.Exec(ctx, q, username); err != nil {
		if res.RowsAffected() == 0 {
			return core.ErrUserNotFound
		}

		return err
	}

	return nil
}

func (r *UsersRepo) Create(ctx context.Context, user *core.User) error {
	q := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"

	err := r.db.QueryRow(ctx, q, user.Username, user.Password).Scan(&user.ID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" && pgErr.ConstraintName == "users_username_key" {
			return core.ErrUserAlreadyExists
		}
	}

	return err
}

func (r *UsersRepo) GetByCredentials(ctx context.Context, username, password string) (core.User, error) {
	q := "SELECT id, username, password FROM users WHERE username=$1 and password=$2"

	var user core.User
	if err := r.db.QueryRow(ctx, q, username, password).Scan(&user.ID, &user.Username, &user.Password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core.User{}, core.ErrUserNotFound
		}

		return core.User{}, err
	}

	return user, nil
}
