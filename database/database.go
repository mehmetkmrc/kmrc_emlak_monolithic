package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DBPool *pgxpool.Pool
var DBErr error

func InitiliazeDatabaseConnection() {
	// Load env file (ignore error if running in production)
	_ = godotenv.Load()

	// Read env variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Check required values
	if dbUser == "" || dbPass == "" || dbHost == "" || dbPort == "" || dbName == "" {
		fmt.Fprintf(os.Stderr, "❌ Database env variables are not fully set\n")
		os.Exit(1)
	}

	// Build Postgres connection URL
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	// Init connection pool
	DBPool, DBErr = pgxpool.New(context.Background(), url)
	if DBErr != nil {
		fmt.Fprintf(os.Stderr, "❌ Unable to create PGX pool: %v\n", DBErr)
		os.Exit(1)
	}

	// Test connection
	err := DBPool.Ping(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database ping failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🚀 Connected to PostgreSQL successfully")
}
