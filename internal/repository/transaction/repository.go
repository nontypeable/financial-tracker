package transaction

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
	"github.com/nontypeable/financial-tracker/internal/domain/transaction"
	apperror "github.com/nontypeable/financial-tracker/internal/errors"
)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) transaction.Repository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) Create(ctx context.Context, transaction *transaction.Transaction) (uuid.UUID, error) {
	query := `
		INSERT INTO transactions (account_id, amount, type, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var id uuid.UUID
	err := r.pool.QueryRow(ctx, query,
		transaction.AccountID,
		transaction.Amount,
		transaction.Type,
		transaction.Description,
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.NotNullViolation, pgerrcode.CheckViolation:
				return uuid.Nil, apperror.ErrInvalidInput
			case pgerrcode.ForeignKeyViolation:
				return uuid.Nil, apperror.ErrAccountNotFound
			}
		}
		return uuid.Nil, fmt.Errorf("create transaction: %w", err)
	}

	return id, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*transaction.Transaction, error) {
	query := `
		SELECT id, account_id, amount, type, description, created_at, updated_at, deleted_at
		FROM transactions
		WHERE id = $1 AND deleted_at IS NULL
	`

	var t transaction.Transaction
	var deletedAt pgtype.Timestamptz

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&t.ID,
		&t.AccountID,
		&t.Amount,
		&t.Type,
		&t.Description,
		&t.CreatedAt,
		&t.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrTransactionNotFound
		}
		return nil, fmt.Errorf("get transaction by id: %w", err)
	}

	if deletedAt.Valid {
		t.DeletedAt = &deletedAt.Time
	}

	return &t, nil
}

func (r *repository) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]*transaction.Transaction, error) {
	query := `
		SELECT id, account_id, amount, type, description, created_at, updated_at, deleted_at
		FROM transactions
		WHERE account_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("get transactions by account_id: %w", err)
	}
	defer rows.Close()

	var transactions []*transaction.Transaction
	for rows.Next() {
		var t transaction.Transaction
		var deletedAt pgtype.Timestamptz

		err = rows.Scan(
			&t.ID,
			&t.AccountID,
			&t.Amount,
			&t.Type,
			&t.Description,
			&t.CreatedAt,
			&t.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan transaction row: %w", err)
		}

		if deletedAt.Valid {
			t.DeletedAt = &deletedAt.Time
		}

		transactions = append(transactions, &t)
	}

	return transactions, nil
}

func (r *repository) Update(ctx context.Context, transaction *transaction.Transaction) error {
	query := `
		UPDATE transactions
		SET amount = $1,
			type = $2,
			description = $3,
			updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL
		RETURNING updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		transaction.Amount,
		transaction.Type,
		transaction.Description,
		transaction.ID,
	).Scan(&transaction.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrTransactionNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.NotNullViolation, pgerrcode.CheckViolation:
				return apperror.ErrInvalidInput
			}
		}

		return fmt.Errorf("update transaction: %w", err)
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, accountID, id uuid.UUID) error {
	query := `
		UPDATE transactions
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND account_id = $2 AND deleted_at IS NULL
	`

	result, err := r.pool.Exec(ctx, query, id, accountID)
	if err != nil {
		return fmt.Errorf("delete transaction: %w", err)
	}

	if result.RowsAffected() == 0 {
		return apperror.ErrTransactionNotFound
	}

	return nil
}
