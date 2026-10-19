package repository

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spearszx/go-wallet/monolith/internal/config"
	"github.com/spearszx/go-wallet/monolith/internal/database"
	"github.com/spearszx/go-wallet/monolith/internal/user/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var (
	testConfig *config.Config
	testPool   *pgxpool.Pool
	testRepo   UserRepository
)

var (
	testFirstNames = []string{
		"ivan", "maria", "alexey", "olga", "dmitry",
		"ekaterina", "sergey", "anna", "nikolay", "tatyana",
	}

	testLastNames = []string{
		"ivanov", "petrova", "smirnov", "kuznetsova", "popov",
		"sokolova", "lebedev", "novikova", "morozov", "volkova",
	}

	eagerDeleteQuery = fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, postgresUsersTableName)
)

func TestMain(m *testing.M) {
	testConfig = config.LoadConfig()

	pool, err := database.NewWithRetries(testConfig.DbOptions)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	testPool = pool
	testRepo = NewPostgresqlUserRepository(pool)

	code := m.Run()

	testPool.Close()
	os.Exit(code)

}

func generateTestUser() (*model.User, error) {
	firstName := testFirstNames[rand.Intn(len(testFirstNames))]
	lastName := testLastNames[rand.Intn(len(testLastNames))]
	email := fmt.Sprintf("%s.%s%d@example.com", firstName, lastName, rand.Int())

	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error to generate hash from password: %v", err)
	}

	user := &model.User{
		ID:           uuid.New().String(),
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: string(hash),
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		DeletedAt:    nil,
	}

	return user, nil
}

func prepareTestUser(t *testing.T) *model.User {
	t.Helper()

	user, err := generateTestUser()
	if err != nil {
		t.Fatalf("generate user: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	created, err := testRepo.Create(ctx, user)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
		defer cancel()

		testPool.Exec(ctx, eagerDeleteQuery, created.ID)
	})

	return created
}

func TestCreateUser(t *testing.T) {
	user := &model.User{
		ID:           "555D9AF86CA7422E9ED7C103A4A4881D",
		FirstName:    "Ivan",
		LastName:     "Ivanov",
		Email:        "ivan.ivanov1@example.com",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
		IsActive:     true,
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
		defer cancel()

		testPool.Exec(ctx, eagerDeleteQuery, user.ID)
	})

	cases := []struct {
		name           string
		input          *model.User
		expectedResult error
	}{
		{
			"normal_user_data",
			user,
			nil,
		},
		{
			"user_already_exists",
			user,
			ErrUserAlreadyExists,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
			defer cancel()

			_, err := testRepo.Create(ctx, tCase.input)
			require.ErrorIs(t, err, tCase.expectedResult)
		})
	}
}

func TestCreateUserWithDBIntegration(t *testing.T) {
	t.Parallel()

	user, err := generateTestUser()
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(t.Context(), time.Second*4)
	defer cancel()

	// create new user
	created, err := testRepo.Create(ctx, user)

	require.NoError(t, err)
	require.NotNil(t, created)

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(t.Context(), time.Second*4)
		defer cancel()

		testPool.Exec(ctx, eagerDeleteQuery, created.ID)
	})

	require.NotNil(t, created.CreatedAt)
	require.NotNil(t, created.CreatedAt)
}

func TestFindByID(t *testing.T) {

	created := prepareTestUser(t)

	cases := []struct {
		name           string
		input          string
		expectedResult error
	}{
		{
			"user_will_be_found_by_id",
			"000000000000000000000000000000000",
			ErrUserNotFound,
		},
		{
			"user_won't_be_found_by_id",
			created.ID,
			nil,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
			defer cancel()

			_, err := testRepo.FindByID(ctx, tCase.input)
			require.ErrorIs(t, err, tCase.expectedResult)
		})
	}
}

func TestFindByIDWithDBIntegration(t *testing.T) {
	t.Parallel()

	created := prepareTestUser(t)

	ctx, cancel := context.WithTimeout(t.Context(), time.Second*4)
	defer cancel()

	found, err := testRepo.FindByID(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, found.ID, created.ID)
}

func TestUpdate(t *testing.T) {
	uninserted, err := generateTestUser()
	require.NoError(t, err)

	created := prepareTestUser(t)
	updated := &model.User{
		ID:           created.ID,
		FirstName:    created.LastName,
		LastName:     created.FirstName,
		Email:        created.Email,
		PasswordHash: created.PasswordHash,
		IsActive:     created.IsActive,
		CreatedAt:    created.CreatedAt,
		UpdatedAt:    created.UpdatedAt,
		DeletedAt:    created.DeletedAt,
	}

	cases := []struct {
		name           string
		input          *model.User
		expectedResult error
	}{
		{
			"user_won't_be_found_to_update",
			uninserted,
			ErrUserNotFound,
		},
		{
			"user_will_be_updated_successfully",
			updated,
			nil,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
			defer cancel()

			err := testRepo.Update(ctx, tCase.input)
			require.ErrorIs(t, err, tCase.expectedResult)
		})
	}
}

func TestUpdateWithDbIntegration(t *testing.T) {
	t.Parallel()

	created := prepareTestUser(t)

	oldFirstName, oldLastName := created.FirstName, created.LastName
	created.FirstName, created.LastName = created.LastName, created.FirstName

	ctx, cancel := context.WithTimeout(t.Context(), time.Second*4)
	defer cancel()

	err := testRepo.Update(ctx, created)
	require.NoError(t, err)

	found, err := testRepo.FindByID(ctx, created.ID)
	require.NoError(t, err)

	assert.Equal(t, found.ID, created.ID)

	require.NotEqual(t, oldFirstName, found.FirstName)
	require.NotEqual(t, oldLastName, found.LastName)
}

func TestFindByEmail(t *testing.T) {
	created := prepareTestUser(t)

	cases := []struct {
		name           string
		input          string
		expectedResult error
	}{
		{
			"user_will_be_found_by_email",
			"invalid_email@example.com",
			ErrUserNotFound,
		},
		{
			"user_won't_be_found_by_email",
			created.Email,
			nil,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
			defer cancel()

			_, err := testRepo.FindByEmail(ctx, tCase.input)
			require.ErrorIs(t, err, tCase.expectedResult)
		})
	}
}

func TestFindByEmailWithDBIntegration(t *testing.T) {
	t.Parallel()

	created := prepareTestUser(t)

	ctx, cancel := context.WithTimeout(t.Context(), time.Second*4)
	defer cancel()

	found, err := testRepo.FindByEmail(ctx, created.Email)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, found.ID, created.ID)
}

func TestDeleteWithIntegration(t *testing.T) {
	t.Parallel()

	created := prepareTestUser(t)

	ctx, cancel := context.WithTimeout(t.Context(), time.Second*4)
	defer cancel()

	err := testRepo.Delete(ctx, created.ID)
	require.NoError(t, err)

	found, err := testRepo.FindByID(ctx, created.ID)
	require.NoError(t, err)

	require.NotNil(t, found)
	require.Equal(t, created.ID, found.ID)
	require.NotNil(t, found.DeletedAt)
}

func TestIsDuplicateKeyError(t *testing.T) {
	cases := []struct {
		name           string
		input          error
		expectedResult bool
	}{
		{
			"is_duplicate_key_will_return_true",
			&pgconn.PgError{
				Code: pgUniqueViolation,
			},
			true,
		},
		{
			"is_duplicate_key_will_return_false",
			nil,
			false,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			result := isDuplicateKeyError(tCase.input)
			require.Equal(t, result, tCase.expectedResult)
		})
	}
}
