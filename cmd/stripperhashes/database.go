package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

// initDatabase initializes the PostgreSQL database
func initDatabase() (*pgx.Conn, error) {
	adminStrBytes, err := os.ReadFile("./adminConn.txt")
	if err != nil {
		return nil, fmt.Errorf("error reading connection string file: %v", err)
	}

	adminStr := string(adminStrBytes)

	connStrBytes, err := os.ReadFile("./connstring.txt")
	if err != nil {
		return nil, fmt.Errorf("error reading connection string file: %v", err)
	}

	connStr := string(connStrBytes)

	adminConn, err := pgx.Connect(context.Background(), adminStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer adminConn.Close(context.Background())

	// Create the database if it does not exist
	createDatabaseSQL := `CREATE DATABASE libergodb;`
	_, err = adminConn.Exec(context.Background(), createDatabaseSQL)
	if err != nil {
		return nil, fmt.Errorf("error creating database: %v", err)
	}

	// Connect to the newly created database
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Create the table in the public schema if it does not exist
	createTableSQL := `CREATE TABLE public.permutations (
        id uuid PRIMARY KEY,
        startArray TEXT,
        endArray TEXT,
        packageName TEXT,
        permName TEXT,
        reportedToAPI BOOLEAN,
        processed BOOLEAN,
        arrayLength INTEGER,
        numberOfPermutations INTEGER DEFAULT 0
    );`

	_, err = conn.Exec(context.Background(), createTableSQL)
	if err != nil {
		err := conn.Close(context.Background())
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	return conn, nil
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

// getByteArrayRanges retrieves the unprocessed byte array ranges from the database
// getByteArrayRange retrieves all rows where numberOfPermutations is equal to 1
func getByteArrayRanges(db *pgx.Conn) ([]struct {
	ID                   string
	StartArray           []byte
	EndArray             []byte
	NumberOfPermutations int
	ArrayLength          int
}, error) {
	rows, err := db.Query(context.Background(), "SELECT id, startArray, endArray, numberOfPermutations, arrayLength FROM public.permutations WHERE numberOfPermutations = 1;")
	if err != nil {
		return nil, fmt.Errorf("error querying rows: %v", err)
	}
	defer rows.Close()

	var results []struct {
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

		results = append(results, struct {
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
	defer conn.Close(context.Background())

	var count int64
	err = conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM public.permutations WHERE numberOfPermutations = 1;").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error querying count: %v", err)
	}

	return count, nil
}
