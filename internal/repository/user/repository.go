package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nontypeable/financial-tracker/internal/domain/user"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) user.Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *user.User) (uuid.UUID, error) {
	query := `
        INSERT INTO users (email, password_hash, first_name, last_name)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	var id uuid.UUID

	err := r.db.QueryRowContext(ctx, query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
	).Scan(&id)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	query := `
        SELECT id, email, password_hash, first_name, last_name, created_at, updated_at, deleted_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `
	var u user.User
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.FirstName,
		&u.LastName,
		&u.CreatedAt,
		&u.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	if deletedAt.Valid {
		u.DeletedAt = &deletedAt.Time
	}

	return &u, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
        SELECT id, email, password_hash, first_name, last_name, created_at, updated_at, deleted_at
        FROM users
        WHERE email = $1 AND deleted_at IS NULL
    `
	var u user.User
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.FirstName,
		&u.LastName,
		&u.CreatedAt,
		&u.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if deletedAt.Valid {
		u.DeletedAt = &deletedAt.Time
	}

	return &u, nil
}

func (r *repository) Update(ctx context.Context, u *user.User) error {
	query := `
        UPDATE users
        SET email = $1,
            password_hash = $2,
            first_name = $3,
            last_name = $4,
            updated_at = NOW()
        WHERE id = $5 AND deleted_at IS NULL
        RETURNING updated_at
    `
	err := r.db.QueryRowContext(ctx, query,
		u.Email,
		u.PasswordHash,
		u.FirstName,
		u.LastName,
		u.ID,
	).Scan(&u.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
        UPDATE users
        SET deleted_at = NOW(), updated_at = NOW()
        WHERE id = $1 AND deleted_at IS NULL
    `
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete user rows affected: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
