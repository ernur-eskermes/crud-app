package psql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"

	"github.com/ernur-eskermes/crud-app/internal/core"
	"github.com/ernur-eskermes/crud-app/pkg/database/postgresql"
	"github.com/google/uuid"
)

type SessionsRepo struct {
	db postgresql.Client
}

func NewSessionsRepo(db postgresql.Client) *SessionsRepo {
	return &SessionsRepo{
		db: db,
	}
}

func (r *SessionsRepo) Create(ctx context.Context, session core.RefreshSession) error {
	q := "INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3) RETURNING id"

	return r.db.QueryRow(ctx, q, session.UserID, session.Token, session.ExpiresAt).Scan(&session.ID)
}

func (r *SessionsRepo) GetByToken(ctx context.Context, token string) (core.RefreshSession, error) {
	q := "SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token=$1"

	var t core.RefreshSession
	if err := r.db.QueryRow(ctx, q, token).Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core.RefreshSession{}, core.ErrTokenNotFound
		}

		return t, err
	}

	return t, nil
}

func (r *SessionsRepo) Delete(ctx context.Context, userID uuid.UUID) error {
	q := "DELETE FROM refresh_tokens WHERE user_id=$1"
	_, err := r.db.Exec(ctx, q, userID)

	return err
}
