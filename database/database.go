package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool
var DBErr error
var url = "postgres://postgres:KMRC-computer-2025@localhost:5432/kmrc_emlak?sslmode=disable"

func InitiliazeDatabaseConnection() {

	DBPool, DBErr = pgxpool.New(context.Background(), url)
	if DBErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", DBErr)
		os.Exit(1)
	}
}
