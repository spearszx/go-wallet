package main

import (
	"fmt"

	"github.com/spearszx/go-wallet/monolith/internal/config"
	"github.com/spearszx/go-wallet/monolith/internal/database"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Println(cfg.ConnectionString())

	pool, err := database.NewWithRetries(cfg.DbOptions)
	if err != nil {
		fmt.Println(err)
	}
	defer pool.Close()
	

	fmt.Println("end")

	_ = pool
}
