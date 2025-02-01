package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

type Permutation struct {
	ID                   string `json:"id"`
	StartArray           []byte `json:"start_array"`
	EndArray             []byte `json:"end_array"`
	PackageName          string `json:"package_name"`
	PermName             string `json:"perm_name"`
	ReportedToAPI        bool   `json:"reported_to_api"`
	Processed            bool   `json:"processed"`
	ArrayLength          int    `json:"array_length"`
	NumberOfPermutations int64  `json:"number_of_permutations"`
}

// initDatabase initializes the PostgreSQL database
func initConnection() (*pgx.Conn, error) {
	connStrBytes, err := os.ReadFile("./connstring.txt")
	if err != nil {
		return nil, fmt.Errorf("error reading connection string file: %v", err)
	}

	connStr := string(connStrBytes)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return conn, nil
}

// closeConnection closes the PostgreSQL database connection
func closeConnection(db *pgx.Conn) error {
	err := db.Close(context.Background())
	if err != nil {
		return fmt.Errorf("error closing connection: %v", err)
	}
	return nil
}

// getByteArrayRange retrieves the unprocessed byte array ranges from the database
func getByteArrayRange(db *pgx.Conn) (*Permutation, error) {
	row := db.QueryRow(context.Background(), "SELECT id, startArray, endArray, numberOfPermutations, arrayLength FROM public.permutations LIMIT 1;")

	var p Permutation
	var startArrayStr, endArrayStr string
	if err := row.Scan(&p.ID, &startArrayStr, &endArrayStr, &p.NumberOfPermutations, &p.ArrayLength); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No more rows to process
		}
		return nil, fmt.Errorf("error scanning row: %v", err)
	}

	var err error
	p.StartArray, err = convertToByteArray(startArrayStr)
	if err != nil {
		return nil, fmt.Errorf("error converting start array: %v", err)
	}

	p.EndArray, err = convertToByteArray(endArrayStr)
	if err != nil {
		return nil, fmt.Errorf("error converting end array: %v", err)
	}

	return &p, nil
}

// removeItem marks a row as processed in the database
func removeItem(db *pgx.Conn, id string) error {
	_, err := db.Exec(context.Background(), "DELETE FROM permutations WHERE id = $1;", id)
	if err != nil {
		return fmt.Errorf("error marking row as processed: %v", err)
	}

	return nil
}

// removeProcessedRows removes the processed rows from the database and compacts it
func removeProcessedRows(db *pgx.Conn) error {
	_, err := db.Exec(context.Background(), "DELETE FROM permutations WHERE processed = true;")
	if err != nil {
		return fmt.Errorf("error deleting processed rows: %v", err)
	}

	fmt.Println("Processed rows removed.")
	return nil
}
