package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func initDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "../permutations.db")
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	return db, nil
}

func getByteArrayRanges(db *sql.DB) ([]struct {
	ID                   string
	StartArray           []byte
	EndArray             []byte
	NumberOfPermutations int
	ArrayLength          int
}, error) {
	rows, err := db.Query("SELECT id, startArray, endArray, numberOfPermutations, arrayLength FROM permutations WHERE processed = 0")
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	var ranges []struct {
		ID                   string
		StartArray           []byte
		EndArray             []byte
		NumberOfPermutations int
		ArrayLength          int
	}
	for rows.Next() {
		var id, startArrayStr, endArrayStr string
		var numberOfPermutations, arrayLength int
		if err := rows.Scan(&id, &startArrayStr, &endArrayStr, &numberOfPermutations, &arrayLength); err != nil {
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

		ranges = append(ranges, struct {
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
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return ranges, nil
}

func markAsProcessed(db *sql.DB, id string) error {
	_, err := db.Exec("UPDATE permutations SET processed = 1 WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error marking row as processed: %v", err)
	}
	fmt.Printf("Row with ID %s marked as processed.\n", id)
	return nil
}

func countUnprocessedRows(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM permutations WHERE processed = 0").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting unprocessed rows: %v", err)
	}
	return count, nil
}
