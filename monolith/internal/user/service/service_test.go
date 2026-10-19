package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spearszx/go-wallet/monolith/internal/config"
	"github.com/spearszx/go-wallet/monolith/internal/user/repository"
)

var (
	testConfig *config.Config
)

func TestMain(m *testing.M) {
	testConfig = config.LoadConfig()
	code := m.Run()

	os.Exit(code)
}

func createUserRepository(t *testing.T) repository.UserRepository {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	pool, err := pgxpool.New(ctx, testConfig.ConnectionString())
	if err != nil {
		t.Error(err)
	}

	return repository.NewPostgresqlUserRepository(pool)
}

// func TestCreateUser(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		in   model.CreateUserRequest
// 		out  error
// 	}{
// 		{
// 			name: "",
// 		},
// 	}

// 	userRepo := createUserRepository(t)
// 	userSvc := NewUserSerivce(userRepo)

// 	userSvc.Create()

// }
