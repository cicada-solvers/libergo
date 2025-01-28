package main

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var dbMutex sync.Mutex

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
		packageFileName TEXT,
		permFileName TEXT,
		reportedToAPI BOOLEAN,
		processed BOOLEAN,
		arrayLength INTEGER
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	return db, nil
}

// insertWithRetry inserts a record into the database with retry logic
func insertWithRetry(db *sql.DB, query string, args ...interface{}) error {
	const maxRetries = 100
	const retryDelay = time.Second

	for i := 0; i < maxRetries; i++ {
		dbMutex.Lock()
		_, err := db.Exec(query, args...)
		dbMutex.Unlock()
		if err == nil {
			return nil
		}
		if strings.Contains(err.Error(), "database is locked") {
			time.Sleep(retryDelay)
			continue
		}
		return err
	}
	return fmt.Errorf("max retries reached, could not insert record")
}

// compactDatabase compacts the SQLite database to reclaim unused space
func compactDatabase(db *sql.DB) error {
	_, err := db.Exec("VACUUM")
	if err != nil {
		return fmt.Errorf("error compacting database: %v", err)
	}
	return nil
}
