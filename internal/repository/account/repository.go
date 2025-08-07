package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nontypeable/financial-tracker/internal/domain/account"
	apperror "github.com/nontypeable/financial-tracker/internal/errors"
)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) account.Repository {
	return &repository{pool: pool}
}

func (r *repository) Create(ctx context.Context, account *account.Account) error {
	query := `
		INSERT INTO accounts (user_id, name, balance)
		VALUES ($1, $2, $3);
	`

	err := r.pool.QueryRow(ctx, query,
		account.ID,
		account.UserID,
		account.Name,
		account.Balance,
	).Scan(&account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.NotNullViolation, pgerrcode.CheckViolation:
				return apperror.ErrInvalidInput
			}
		}
		return fmt.Errorf("failed to create account: %w", err)
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*account.Account, error) {
	query := `
		SELECT id, user_id, name, balance, created_at, updated_at, deleted_at
		FROM accounts
		WHERE id = $1 AND deleted_at IS NULL
	`

	var a account.Account

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&a.ID,
		&a.UserID,
		&a.Name,
		&a.Balance,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrAccountNotFound
		}
		return nil, fmt.Errorf("get account by id: %w", err)
	}

	return &a, nil
}

func (r *repository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*account.Account, error) {
	query := `
		SELECT id, user_id, name, balance, created_at, updated_at, deleted_at
		FROM accounts
		WHERE user_id = $1 AND deleted_at IS NULL
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts by user id: %w", err)
	}
	defer rows.Close()

	var accounts []*account.Account

	for rows.Next() {
		var a account.Account
		err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.Name,
			&a.Balance,
			&a.CreatedAt,
			&a.UpdatedAt,
			&a.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account row: %w", err)
		}
		accounts = append(accounts, &a)
	}

	return accounts, nil
}

func (r *repository) Update(ctx context.Context, account *account.Account) error {
	query := `
		UPDATE accounts
		SET name = $1,
		    balance = $2,
		    updated_at = NOW()
		WHERE id = $3 AND deleted_at IS NULL
		RETURNING updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		account.Name,
		account.Balance,
		account.ID,
	).Scan(&account.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrInvalidInput
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.NotNullViolation, pgerrcode.CheckViolation:
				return apperror.ErrInvalidInput
			}
		}

		return fmt.Errorf("update account: %w", err)
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	query := `
		UPDATE accounts
		SET deleted_at = NOW()
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`

	ct, err := r.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to soft-delete account: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return errors.New("no account found to delete")
	}

	return nil
}
