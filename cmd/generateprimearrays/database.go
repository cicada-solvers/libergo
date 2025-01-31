package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
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

func initConnection() (*pgx.Conn, error) {
	connStrBytes, err := os.ReadFile("./connstring.txt")
	if err != nil {
		return nil, fmt.Errorf("error reading connection string file: %v", err)
	}

	connStr := string(connStrBytes)

	// Connect to the newly created database
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn, nil
}

// insertWithRetry inserts a record into the database with retry logic
func insertWithRetry(db *pgx.Conn, query string) error {
	var err error
	_, err = db.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("error inserting record: %v", err)
	}
	return nil
}
