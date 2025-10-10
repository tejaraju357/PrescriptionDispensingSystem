package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func DBConnect() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	connStr := os.Getenv("DBDSN")
	if connStr == "" {
		log.Fatal("DBDSN environment variable not set")
	}

	// Parse the config
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Failed to parse DB config: %v", err)
	}

	// ✅ Disable pgx statement cache to avoid "stmtcache" errors
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// Connect using the config
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("❌ Unable to connect to database: %v", err)
	}

	DB = pool
	fmt.Println("✅ DB connected successfully (statement caching disabled)")
}
