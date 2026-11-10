package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite" // or modernc.org/sqlite (CGO-free)
)

func InitDB() *sql.DB {
	// Points directly to a file inside your project directory
	db, err := sql.Open("sqlite", "./metro.db?_journal_mode=WAL")
	if err != nil {
		log.Fatalf("Failed to connect to SQLite: %v", err)
	}

	// Limit connections for SQLite concurrent safety
	db.SetMaxOpenConns(10)

	// Create tables if they don't exist yet
	query := `
	CREATE TABLE IF NOT EXISTS metro_lines (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    line_code TEXT UNIQUE NOT NULL,
    line_color_verbose TEXT,
    line_color_family TEXT,
    primary_color_code TEXT,
    secondary_color_code TEXT,
    station_count INTEGER,
    start_station JSON,
    end_station JSON,
    junction_count INTEGER,
    status TEXT
	);

	CREATE TABLE IF NOT EXISTS metro_stations (
    id INTEGER PRIMARY KEY,
    station_code TEXT UNIQUE NOT NULL,
    station_name TEXT NOT NULL,
    station_description TEXT,
    station_type TEXT,
    interchange BOOLEAN DEFAULT 0,
    latitude REAL,
    longitude REAL,
    mobile TEXT,
    landline TEXT,
    opening_time TEXT,
    closing_time TEXT,
    metro_lines JSON,
    prev_next_stations JSON,
    station_facility JSON,
    gates JSON,
    lifts JSON,
    platforms JSON,
    stations_facilities JSON,
    parkings JSON
	);

	CREATE TABLE IF NOT EXISTS line_updates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    stations_updates TEXT,
    line_updates JSON,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	fmt.Println("Database initialized successfully!")
	return db
}

func main() {
	db := InitDB()
	defer db.Close()
}
