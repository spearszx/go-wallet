package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spearszx/go-wallet/monolith/internal/wallet/model"
)

const (
	postgresWalletsTableName = "users"

	pgUniqueViolation = "23505"

	allWalletsFields = "id, user_id, balance, currenct, status, version, created_at, updated_at, deleted_at"
)

var (
	ErrWalletAlreadyExists = errors.New("wallet already exists")
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrInternalFailure     = errors.New("invalid failure")
)

type WalletRepository interface {
	Create(ctx context.Context, user *model.Wallet) (*model.Wallet, error)
	FindByID(ctx context.Context, id string) (*model.Wallet, error)
	FindByUserID(ctx context.Context, userID string) (*model.Wallet, error)
	Update(ctx context.Context, wallet *model.Wallet) error
	Delete(ctx context.Context, id string) error
}

type postgresWalletRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresqlWalletRepository(pool *pgxpool.Pool) WalletRepository {
	return &postgresWalletRepository{
		pool: pool,
	}
}

func (r *postgresWalletRepository) Create(ctx context.Context, wallet *model.Wallet) (*model.Wallet, error) {
	const op = "postgresWalletRepository.Create"

	query := fmt.Sprintf(`
		INSERT INTO %s (id, user_id, balance, currency, status, version)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`, postgresWalletsTableName)

	err := r.pool.QueryRow(ctx, query, wallet.ID, wallet.UserID, wallet.Balance, wallet.Currency, wallet.Status, wallet.Version).
		Scan(&wallet.CreatedAt, &wallet.UpdatedAt)

	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, fmt.Errorf("%s: %w", op, ErrWalletAlreadyExists)
		}
		return nil, fmt.Errorf("%s: %w", op, ErrInternalFailure)
	}

	return wallet, nil
}

func (r *postgresWalletRepository) FindByID(ctx context.Context, id string) (*model.Wallet, error) {
	const op = "postgresWalletRepository.FindByID"

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE id = $1;
	`, allWalletsFields, postgresWalletsTableName)

	user, err := r.scanIntoWallet(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrWalletNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, ErrInternalFailure)
	}

	return user, nil
}

func (r *postgresWalletRepository) FindByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	const op = "postgresWalletRepository.FindByUserID"

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE id = $1;
	`, allWalletsFields, postgresWalletsTableName)

	user, err := r.scanIntoWallet(r.pool.QueryRow(ctx, query, userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrWalletNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, ErrInternalFailure)
	}

	return user, nil

}

func (r *postgresWalletRepository) Update(ctx context.Context, wallet *model.Wallet) error {
	const op = "postgresWalletRepository.Update"

	query := fmt.Sprintf(`
		UPDATE %s
		SET 
			user_id=$1, 
			balance=$2,
			currency=$3,
			status=$4,
			version=$5,
			updated_at=NOW()
		WHERE id = $6;
	`, postgresWalletsTableName)

	tag, err := r.pool.Exec(ctx, query,
		wallet.UserID,
		wallet.Balance,
		wallet.Currency,
		wallet.Status,
		wallet.Version,
		wallet.ID,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, ErrInternalFailure)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrWalletNotFound)
	}

	return nil
}

func (r *postgresWalletRepository) Delete(ctx context.Context, id string) error {
	const op = "postgresWalletRepository.Delete"

	query := fmt.Sprintf(`
		UPDATE %s
		SET 
			deleted_at=NOW()
		WHERE id = $1;
	`, postgresWalletsTableName)

	if _, err := r.pool.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("%s: %w", op, ErrInternalFailure)
	}

	return nil

}

func (r *postgresWalletRepository) scanIntoWallet(row pgx.Row) (*model.Wallet, error) {
	var (
		err    error
		wallet model.Wallet
	)

	err = row.Scan(
		&wallet.ID, &wallet.UserID, &wallet.Balance, &wallet.Currency, &wallet.Status,
		&wallet.Version, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.DeletedAt,
	)

	return &wallet, err
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == pgUniqueViolation // unique_violation
	}
	return false
}
