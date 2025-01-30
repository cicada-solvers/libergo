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
	return db, nil
}

// getByteArrayRanges retrieves the unprocessed byte array ranges from the database
func getByteArrayRange(db *sql.DB) (*struct {
	ID                   string
	StartArray           []byte
	EndArray             []byte
	NumberOfPermutations int
	ArrayLength          int
}, error) {
	row := db.QueryRow("SELECT id, startArray, endArray, numberOfPermutations, arrayLength FROM permutations WHERE processed = 0 LIMIT 1")

	var id, startArrayStr, endArrayStr string
	var numberOfPermutations, arrayLength int
	if err := row.Scan(&id, &startArrayStr, &endArrayStr, &numberOfPermutations, &arrayLength); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No more rows
		}
		return nil, fmt.Errorf("error scanning row: %v", err)
	}

	startArray, err := convertToByteArray(startArrayStr)
	if err != nil {
		return nil, fmt.Errorf("error converting start array: %v", err)
	}

	endArray, err := convertToByteArray(endArrayStr)
	if err != nil {
		return nil, fmt.Errorf("error converting end array: %v", err)
	}

	return &struct {
		ID                   string
		StartArray           []byte
		EndArray             []byte
		NumberOfPermutations int
		ArrayLength          int
	}{
		ID:                   id,
		StartArray:           startArray,
		EndArray:             endArray,
		NumberOfPermutations: numberOfPermutations,
		ArrayLength:          arrayLength,
	}, nil
}

// markAsProcessed marks a row as processed in the database
func markAsProcessed(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM permutations WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error marking row as processed: %v", err)
	}
	fmt.Printf("Row with ID %s marked as processed.\n", id)

	return nil
}

// countUnprocessedRows counts the number of unprocessed rows in the database
func countUnprocessedRows(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM permutations WHERE processed = 0").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting unprocessed rows: %v", err)
	}
	return count, nil
}

// removeProcessedRows removes the processed rows from the database and compacts it
func removeProcessedRows(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM permutations WHERE processed = 1")
	if err != nil {
		return fmt.Errorf("error deleting processed rows: %v", err)
	}

	err = compactDatabase(db)
	if err != nil {
		return fmt.Errorf("error compacting database: %v", err)
	}

	fmt.Println("Processed rows removed and database compacted.")
	return nil
}

// compactDatabase compacts the SQLite database to reclaim unused space
func compactDatabase(db *sql.DB) error {
	_, err := db.Exec("VACUUM")
	if err != nil {
		return fmt.Errorf("error compacting database: %v", err)
	}
	return nil
}
