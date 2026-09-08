package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the connection to Supabase PostgreSQL and creates tables if they don't exist.
func InitDB() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("WARNING: DATABASE_URL environment variable is empty. Database functions will fail.")
		return
	}

	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	// Verify database connection
	err = DB.Ping()
	if err != nil {
		log.Printf("WARNING: Error pinging database: %v. Database functions will fail.", err)
		return
	}
	log.Println("Successfully connected to Supabase PostgreSQL database!")

	// Create tables
	createTables()
}

func createTables() {
	// 1. Create Users Table
	userTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		display_name VARCHAR(255),
		avatar_url TEXT,
		bio TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(userTableSQL)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}
	log.Println("Users table checked/created.")

	// 2. Create Watchlist Table
	watchlistTableSQL := `
	CREATE TABLE IF NOT EXISTS watchlist (
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		anime_id INT NOT NULL,
		title VARCHAR(255) NOT NULL,
		image_url TEXT NOT NULL,
		added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, anime_id)
	);`

	_, err = DB.Exec(watchlistTableSQL)
	if err != nil {
		log.Fatalf("Error creating watchlist table: %v", err)
	}
	log.Println("Watchlist table checked/created.")
}
