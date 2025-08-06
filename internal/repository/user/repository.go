package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nontypeable/financial-tracker/internal/domain/user"
	apperror "github.com/nontypeable/financial-tracker/internal/errors"
)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) user.Repository {
	return &repository{pool: pool}
}

func (r *repository) Create(ctx context.Context, user *user.User) (uuid.UUID, error) {
	query := `
        INSERT INTO users (email, password_hash, first_name, last_name)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	var id uuid.UUID

	err := r.pool.QueryRow(ctx, query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return uuid.Nil, apperror.ErrUserAlreadyExists
			case pgerrcode.NotNullViolation, pgerrcode.CheckViolation:
				return uuid.Nil, apperror.ErrInvalidInput
			default:
				return uuid.Nil, fmt.Errorf("create user: %w", err)
			}
		}
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
	var deletedAt pgtype.Timestamptz

	err := r.pool.QueryRow(ctx, query, id).Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrUserNotFound
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
	var deletedAt pgtype.Timestamptz

	err := r.pool.QueryRow(ctx, query, email).Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrUserNotFound
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

	err := r.pool.QueryRow(ctx, query,
		u.Email,
		u.PasswordHash,
		u.FirstName,
		u.LastName,
		u.ID,
	).Scan(&u.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrUserNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return apperror.ErrUserAlreadyExists
			case pgerrcode.NotNullViolation, pgerrcode.CheckViolation:
				return apperror.ErrInvalidInput
			}
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

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *repository) EmailExists(ctx context.Context, email string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL);`

	var exists bool
	err := r.pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check email existence: %w", err)
	}

	return exists, nil
}
