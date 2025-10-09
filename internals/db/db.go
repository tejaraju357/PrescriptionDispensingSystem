package db

import (
	"context"
	"fmt"

	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func DBConnect() {
	godotenv.Load()

	connStr := os.Getenv("DBDSN")
    if connStr == "" {
        log.Fatal("DATABASE_URL environment variable not set")
    }

	
    pool, err := pgxpool.New(context.Background(), connStr)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }
    DB = pool

    fmt.Println("DB connected successfully")
}