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

// getByteArrayRanges retrieves the unprocessed byte array ranges from the database
func getByteArrayRanges(db *pgx.Conn) ([]Permutation, error) {
	rows, err := db.Query(context.Background(), "SELECT id, startArray, endArray, numberOfPermutations, arrayLength FROM public.permutations WHERE numberOfPermutations = 1;")
	if err != nil {
		return nil, fmt.Errorf("error querying rows: %v", err)
	}
	defer rows.Close()

	var results []Permutation

	for rows.Next() {
		var p Permutation
		var startArrayStr, endArrayStr string
		if err := rows.Scan(&p.ID, &startArrayStr, &endArrayStr, &p.NumberOfPermutations, &p.ArrayLength); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		p.StartArray, err = convertToByteArray(startArrayStr)
		if err != nil {
			return nil, fmt.Errorf("error converting start array: %v", err)
		}

		p.EndArray, err = convertToByteArray(endArrayStr)
		if err != nil {
			return nil, fmt.Errorf("error converting end array: %v", err)
		}

		results = append(results, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return results, nil
}

// removeItem marks a row as processed in the database
func removeItem(db *pgx.Conn, id string) error {
	_, err := db.Exec(context.Background(), "DELETE FROM public.permutations WHERE id = $1;", id)
	if err != nil {
		return fmt.Errorf("error marking row as processed: %v", err)
	}

	return nil
}

// removeProcessedRows removes the processed rows from the database and compacts it
func removeProcessedRows(db *pgx.Conn) error {
	_, err := db.Exec(context.Background(), "DELETE FROM public.permutations WHERE processed = true;")
	if err != nil {
		return fmt.Errorf("error deleting processed rows: %v", err)
	}

	fmt.Println("Processed rows removed.")
	return nil
}

// getCountOfPermutations returns the count of rows where NumberOfPermutations = 1
func getCountOfPermutations() (int64, error) {
	connStrBytes, err := os.ReadFile("./connstring.txt")
	if err != nil {
		return 0, fmt.Errorf("error reading connection string file: %v", err)
	}

	connStr := string(connStrBytes)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return 0, fmt.Errorf("error connecting to database: %v", err)
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			fmt.Printf("error closing connection: %v\n", err)
		}
	}(conn, context.Background())

	var count int64
	err = conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM public.permutations WHERE numberOfPermutations = 1;").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error querying count: %v", err)
	}

	return count, nil
}
