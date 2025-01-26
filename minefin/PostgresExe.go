package minefin

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func (m MediaFileProcessor) PostgreSQL() {
	// Database connection parameters
	host := "localhost"
	port := "5432"
	user := "admin"
	password := "root"
	dbname := "minfin-medialist"

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
