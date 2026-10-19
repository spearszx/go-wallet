package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spearszx/go-wallet/monolith/internal/user/model"
)

const (
	postgresUsersTableName = "users"

	pgUniqueViolation = "23505"

	allUserFields = "id, first_name, last_name, email, password_hash, is_active, created_at, updated_at, deleted_at"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInternalFailure   = errors.New("invalid failure")
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
}

type postgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresqlUserRepository(pool *pgxpool.Pool) UserRepository {
	return &postgresUserRepository{
		pool: pool,
	}
}

func (r *postgresUserRepository) Create(ctx context.Context, u *model.User) (*model.User, error) {
	const op = "postgresUserRepository.Create"

	query := fmt.Sprintf(`
		INSERT INTO %s (id, first_name, last_name, email, password_hash, is_active)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`, postgresUsersTableName)

	err := r.pool.QueryRow(ctx, query, u.ID, u.FirstName, u.LastName, u.Email, u.PasswordHash, u.IsActive).
		Scan(&u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
		}
		return nil, fmt.Errorf("%s: %w", op, ErrInternalFailure)
	}

	return u, nil
}

func (r *postgresUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	const op = "postgresUserRepository.FindByID"

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE id = $1;
	`, allUserFields, postgresUsersTableName)

	user, err := r.scanIntoUser(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, ErrInternalFailure)
	}

	return user, nil
}

func (r *postgresUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	const op = "postgresUserRepository.FindByEmail"

	query := fmt.Sprintf(`
		SELECT %s
		FROM %s
		WHERE email = $1;
	`, allUserFields, postgresUsersTableName)

	user, err := r.scanIntoUser(r.pool.QueryRow(ctx, query, email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, ErrInternalFailure)
	}

	return user, nil

}

func (r *postgresUserRepository) Update(ctx context.Context, u *model.User) error {
	const op = "postgresUserRepository.Update"

	query := fmt.Sprintf(`
		UPDATE %s
		SET 
			first_name=$1, 
			last_name=$2,
			email=$3,
			password_hash=$4,
			is_active=$5,
			updated_at=NOW()
		WHERE id = $6;
	`, postgresUsersTableName)

	tag, err := r.pool.Exec(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.PasswordHash,
		u.IsActive,
		u.ID,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, ErrInternalFailure)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}

	return nil
}

func (r *postgresUserRepository) Delete(ctx context.Context, id string) error {
	const op = "postgresUserRepository.Delete"

	query := fmt.Sprintf(`
		UPDATE %s
		SET 
			deleted_at=NOW()
		WHERE id = $1;
	`, postgresUsersTableName)

	if _, err := r.pool.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("%s: %w", op, ErrInternalFailure)
	}

	return nil
}

func (r *postgresUserRepository) scanIntoUser(row pgx.Row) (*model.User, error) {
	var (
		err  error
		user model.User
	)

	err = row.Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)

	return &user, err
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == pgUniqueViolation // unique_violation
	}
	return false
}
