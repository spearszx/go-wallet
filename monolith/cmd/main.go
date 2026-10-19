package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spearszx/go-wallet/monolith/internal/config"
	"github.com/spearszx/go-wallet/monolith/internal/database"
	"github.com/spearszx/go-wallet/monolith/internal/logger"
	"github.com/spearszx/go-wallet/monolith/internal/server"
	"github.com/spearszx/go-wallet/monolith/internal/user/repository"
	"github.com/spearszx/go-wallet/monolith/internal/user/service"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Println(cfg.ConnectionString())

	logger := logger.LoadLogger(os.Stdout)

	pool, err := database.NewWithRetries(cfg.DbOptions)
	if err != nil {
		fmt.Println(err)
	}
	defer pool.Close()

	// TODO: Передать в server и оттуда в handlers
	validate := validator.New()
	userRepo := repository.NewPostgresqlUserRepository(pool)
	userService := service.NewUserSerivce(userRepo)

	server := server.NewServer(
		server.Options{
			ListenAddr: cfg.ListenAddr,
		},
		logger,
		validate,
		userService,
	)

	log.Fatal(server.ListenAndServe())
}
