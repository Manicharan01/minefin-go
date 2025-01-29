package minefin

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type LibraryManager struct {
	LibraryPath string
}

func (l LibraryManager) connectDB() (*sql.DB, error) {
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
	portInt, _ := strconv.Atoi(port)

	// Create connection URL
	connectionURL := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", user, password, dbname, host, portInt)

	db, err := sql.Open("postgres", connectionURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (l LibraryManager) GetTableNames() []string {
	dbConnector, err := l.connectDB()
	if err != nil {
		log.Fatalf("Error while connecting to db: %v", err)
	}
	defer dbConnector.Close()

	rows, err := dbConnector.Query(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`)
	if err != nil {
		log.Fatalf("Error while getting table names: %v", err)
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatalf("Error while getting table names: %v", err)
		}

		tableNames = append(tableNames, tableName)
	}

	return tableNames
}

func (l LibraryManager) CreateTable(tableName string) {
	dbConnector, err := l.connectDB()
	if err != nil {
		log.Fatalf("Error while connecting to DB: %v", err)
	}
	defer dbConnector.Close()

	if tableName == "users" {
		createNewTableQuery := `
		CREATE TABLE Users (
			user_id UUID PRIMARY KEY,
			username VARCHAR(255),
			password_hash TEXT,
			email VARCHAR(255),
			created_at TIMESTAMP
		);`

		_, err := dbConnector.Exec(createNewTableQuery)
		if err != nil {
			log.Fatalf("Error while creating user table: %v", err)
		}

		fmt.Println("Users table created successfully")
	} else if tableName == "mediaitems" {
		createNewTableQuery := `
		CREATE TABLE MediaItems (
			media_id UUID PRIMARY KEY,
			title VARCHAR(255),
			type VARCHAR(50),
			file_path TEXT,
			duration INTEGER,
			format VARCHAR(50),
			created_at TIMESTAMP,
			metadata JSONB
		);`

		_, err := dbConnector.Exec(createNewTableQuery)
		if err != nil {
			log.Fatalf("Error while creating mediaitems table: %v", err)
		}

		fmt.Println("MediaItems table created successfully")
	} else if tableName == "usermediaprogress" {
		createNewTableQuery := `
		CREATE TABLE UserMediaProgress (
			user_id UUID,
			media_id UUID,
			position INTEGER,
			last_watched TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES Users(user_id),
			FOREIGN KEY (media_id) REFERENCES MediaItems(media_id)
		);`

		_, err := dbConnector.Exec(createNewTableQuery)
		if err != nil {
			log.Fatalf("Error while creating usermediaprogress table: %v", err)
		}

		fmt.Println("UserMediaProgress table created successfully")
	}
}
