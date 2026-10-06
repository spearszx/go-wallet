package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	defaultDbUserValue     = "postgres"
	defaultDbPasswordValue = ""
	defaultDbHostValue     = "localhost"
	defaultDbPortValue     = "5432"
	defaultDbNameValue     = "postgres"
	defaultMaxRetriesValue = 4

	defaultServerPort = ":8080"
)

var (
	isInitialised bool
	config        *Config
)

type Config struct {
	*DbOptions
	*ServerOptions
}

type DbOptions struct {
	DbUser       string
	DbPort       string
	DbPassword   string
	DbHost       string
	DbName       string
	SslMode      bool
	DbMaxRetries int
}

type ServerOptions struct {
	ListenAddr string
}

func (Db *DbOptions) ConnectionString() string {
	// return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%v",
	// 	Db.DbUser, Db.DbPassword, Db.DbHost, Db.DbPort, Db.DbName, Db.SslMode)

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		Db.DbUser, Db.DbPassword, Db.DbHost, Db.DbPort, Db.DbName)

}

func LoadConfig() *Config {
	if !isInitialised {
		initConfig()
	}
	return config
}

func initConfig() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning .env file not found, using environment variables")
	}

	config = &Config{
		DbOptions: &DbOptions{
			DbUser:       getEnv("POSTGRES_USER", defaultDbUserValue),
			DbPassword:   getEnv("POSTGRES_PASSWORD", defaultDbPasswordValue),
			DbHost:       getEnv("POSTGRES_HOST", defaultDbHostValue),
			DbPort:       getEnv("POSTGRES_PORT", defaultDbPortValue),
			DbName:       getEnv("POSTGRES_DB", defaultDbNameValue),
			DbMaxRetries: getIntEnv("DB_MAX_RETRIES", defaultMaxRetriesValue),
		},
		ServerOptions: &ServerOptions{
			ListenAddr: getEnv("SERVER_PORT", defaultServerPort),
		},
	}

	isInitialised = true
}

func getIntEnv(key string, defaultValue int) int {
	stringVal := os.Getenv(key)
	if stringVal == "" {
		return defaultValue
	}

	intVal, err := strconv.Atoi(stringVal)
	if err != nil {
		return defaultValue
	}

	return intVal
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
