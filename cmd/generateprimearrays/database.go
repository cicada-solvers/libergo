package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type Permutation struct {
	ID                   string `json:"id"`
	StartArray           string `json:"start_array"`
	EndArray             string `json:"end_array"`
	PackageName          string `json:"package_name"`
	PermName             string `json:"perm_name"`
	ReportedToAPI        bool   `json:"reported_to_api"`
	Processed            bool   `json:"processed"`
	ArrayLength          int    `json:"array_length"`
	NumberOfPermutations int64  `json:"number_of_permutations"`
}

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
		_, err := fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		if err != nil {
			return nil, err
		}
		os.Exit(1)
	}
	defer func(adminConn *pgx.Conn, ctx context.Context) {
		err := adminConn.Close(ctx)
		if err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Error closing database connection: %v\n", err)
			if err != nil {
				return
			}
		}
	}(adminConn, context.Background())

	// Create the database if it does not exist
	createDatabaseSQL := `CREATE DATABASE libergodb;`
	_, err = adminConn.Exec(context.Background(), createDatabaseSQL)
	if err != nil {
		return nil, fmt.Errorf("error creating database: %v", err)
	}

	// Connect to the newly created database
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		if err != nil {
			return nil, err
		}
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
        arrayLength BIGINT,
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
		_, err := fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		if err != nil {
			return nil, err
		}
		os.Exit(1)
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

// insertWithRetry inserts a record into the database with retry logic
func insertRecord(db *pgx.Conn, perm Permutation) error {
	query := `INSERT INTO public.permutations(
                                id, 
                                startArray, 
                                endArray, 
                                packageName, 
                                permName, 
                                reportedToAPI, 
                                processed, 
                                arrayLength, 
                                numberOfPermutations) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`

	_, err := db.Exec(context.Background(),
		query,
		perm.ID,
		perm.StartArray,
		perm.EndArray,
		perm.PackageName,
		perm.PermName,
		perm.ReportedToAPI,
		perm.Processed,
		perm.ArrayLength,
		perm.NumberOfPermutations)

	if err != nil {
		return fmt.Errorf("error inserting record: %v", err)
	}
	return nil
}
