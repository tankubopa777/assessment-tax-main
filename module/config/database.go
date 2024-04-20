package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func NewDatabaseConnection() (*sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	
	AdminSettingsTable(db)

	return db, nil
}

func AdminSettingsTable(db *sql.DB) {
	tableCreationQuery := `CREATE TABLE IF NOT EXISTS admin_settings(
		id SERIAL PRIMARY KEY,
		personal_deduction NUMERIC DEFAULT 60000,
		k_receipt_limit NUMERIC DEFAULT 50000
	)`

	_, err := db.Exec(tableCreationQuery)
	if err != nil {
		log.Fatalf("Error creating admin_settings table: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM admin_settings").Scan(&count)
	if err != nil {
		log.Fatalf("Error counting admin settings: %v", err)
	}

	if count == 0 {
		_, err = db.Exec("INSERT INTO admin_settings (personal_deduction, k_receipt_limit) VALUES ($1, $2)", 60000, 50000)
		if err != nil {
			log.Fatalf("Error inserting default admin settings: %v", err)
		}
	}
}
