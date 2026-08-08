// Package user implements the Postgres-backed User repository.
package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Repo -.
type Repo struct {
	*postgres.Postgres
}

// New returns a User repository instrumented with OpenTelemetry tracing spans.
func New(pg *postgres.Postgres) repo.UserRepo {
	return newTraced(&Repo{pg})
}

// Store -.
func (r *Repo) Store(ctx context.Context, user *entity.User) error {
	const sql = `INSERT INTO users (id, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4)`

	_, err := r.Pool.Exec(ctx, sql, user.ID, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entity.ErrUserAlreadyExists
		}

		return fmt.Errorf("UserRepo - Store - r.Pool.Exec: %w", err)
	}

	return nil
}

// GetByID -.
func (r *Repo) GetByID(ctx context.Context, id string) (entity.User, error) {
	return r.getUser(ctx, "id", id)
}

// GetByEmail -.
func (r *Repo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	return r.getUser(ctx, "email", email)
}

func (r *Repo) getUser(ctx context.Context, column, value string) (entity.User, error) {
	sql := `SELECT id, email, created_at, updated_at
		FROM users
		WHERE ` + column + ` = $1`

	var user entity.User

	err := r.Pool.QueryRow(ctx, sql, value).
		Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, entity.ErrUserNotFound
		}

		return entity.User{}, fmt.Errorf("UserRepo - getUser - r.Pool.QueryRow: %w", err)
	}

	return user, nil
}
