package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// initDatabase initializes the SQLite database
func initDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:../permutations.db?_journal_mode=WAL&_mutex=full")
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS permutations (
		id TEXT PRIMARY KEY,
		startArray TEXT,
		endArray TEXT,
		packageName TEXT,
		permName TEXT,
		reportedToAPI BOOLEAN,
		processed BOOLEAN,
		arrayLength INTEGER,
		numberOfPermutations INTEGER DEFAULT 0
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	return db, nil
}

// compactDatabase compacts the SQLite database to reclaim unused space
func compactDatabase(db *sql.DB) error {
	_, err := db.Exec("VACUUM")
	if err != nil {
		return fmt.Errorf("error compacting database: %v", err)
	}
	return nil
}
