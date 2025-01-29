package minefin

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func (m MediaFileProcessor) PostgreSQL() {
	// Loading .ENV file
	err := godotenv.Load(".ENV")
	if err != nil {
		log.Fatalf("Error while loading .ENV variables: %v", err)
	}

	// Database connection parameters
	host := os.Getenv("POSTGRES_URL")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB_NAME")

	// Create connection URL
	connectionURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		user, password, host, port, dbname)

	// Create connection pool
	config, err := pgxpool.ParseConfig(connectionURL)
	if err != nil {
		log.Fatalf("Unable to parse database config: %v", err)
	}

	// Establish connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer pool.Close()

	// Test connection
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")

	// Example query
	rows, err := pool.Query(context.Background(), "SELECT version()")
	if err != nil {
		log.Fatalf("Query error: %v", err)
	}
	defer rows.Close()

	var version string
	for rows.Next() {
		err := rows.Scan(&version)
		if err != nil {
			log.Fatalf("Scan error: %v", err)
		}
		fmt.Printf("PostgreSQL Version: %s\n", version)
	}
}
