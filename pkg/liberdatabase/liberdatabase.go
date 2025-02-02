package liberdatabase

import (
	"config"
	"context"
	"conversion"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

// InitDatabase initializes the PostgreSQL database
func InitDatabase() (*pgx.Conn, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}

	adminStr := cfg.AdminConnectionString
	connStr := cfg.GeneralConnectionString

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

	factorError := InitFactor(conn)
	if factorError != nil {
		return nil, err
	}

	return conn, nil
}

func InitFactor(conn *pgx.Conn) error {
	// Create the table in the public schema if it does not exist
	createTableSQL := `CREATE TABLE public.factors (
		id uuid PRIMARY KEY,
		factor TEXT,
		mainid uuid
	);`

	_, err := conn.Exec(context.Background(), createTableSQL)
	if err != nil {
		err := conn.Close(context.Background())
		if err != nil {
			return err
		}
		return fmt.Errorf("error creating table: %v", err)
	}

	return nil
}

// InitConnection initializes the PostgreSQL database
func InitConnection() (*pgx.Conn, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}

	connStr := cfg.GeneralConnectionString

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return conn, nil
}

// CloseConnection closes the PostgreSQL database connection
func CloseConnection(db *pgx.Conn) error {
	err := db.Close(context.Background())
	if err != nil {
		return fmt.Errorf("error closing connection: %v", err)
	}
	return nil
}

// GetByteArrayRange retrieves the unprocessed byte array ranges from the database
func GetByteArrayRange(db *pgx.Conn) (*ReadPermutation, error) {
	row := db.QueryRow(context.Background(), "SELECT id, startArray, endArray, numberOfPermutations, arrayLength FROM public.permutations LIMIT 1;")

	var p ReadPermutation
	var startArrayStr, endArrayStr string
	if err := row.Scan(&p.ID, &startArrayStr, &endArrayStr, &p.NumberOfPermutations, &p.ArrayLength); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No more rows to process
		}
		return nil, fmt.Errorf("error scanning row: %v", err)
	}

	var err error
	p.StartArray, err = conversion.ConvertToByteArray(startArrayStr)
	if err != nil {
		return nil, fmt.Errorf("error converting start array: %v", err)
	}

	p.EndArray, err = conversion.ConvertToByteArray(endArrayStr)
	if err != nil {
		return nil, fmt.Errorf("error converting end array: %v", err)
	}

	return &p, nil
}

// GetByteArrayRanges retrieves the unprocessed byte array ranges from the database
func GetByteArrayRanges(db *pgx.Conn) ([]ReadPermutation, error) {
	rows, err := db.Query(context.Background(), "SELECT id, startArray, endArray, numberOfPermutations, arrayLength FROM public.permutations WHERE numberOfPermutations = 1;")
	if err != nil {
		return nil, fmt.Errorf("error querying rows: %v", err)
	}
	defer rows.Close()

	var results []ReadPermutation

	for rows.Next() {
		var p ReadPermutation
		var startArrayStr, endArrayStr string
		if err := rows.Scan(&p.ID, &startArrayStr, &endArrayStr, &p.NumberOfPermutations, &p.ArrayLength); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		p.StartArray, err = conversion.ConvertToByteArray(startArrayStr)
		if err != nil {
			return nil, fmt.Errorf("error converting start array: %v", err)
		}

		p.EndArray, err = conversion.ConvertToByteArray(endArrayStr)
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

// GetCountOfPermutations returns the count of rows where NumberOfPermutations = 1
func GetCountOfPermutations() (int64, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return 0, fmt.Errorf("error loading config: %v", err)
	}

	connStr := cfg.GeneralConnectionString

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

// GetFactorsByMainID retrieves all factors from the factors table based on the mainid.
func GetFactorsByMainID(db *pgx.Conn, mainId string) ([]Factor, error) {
	query := `SELECT id, factor, mainid FROM public.factors WHERE mainid = $1`
	rows, err := db.Query(context.Background(), query, mainId)
	if err != nil {
		return nil, fmt.Errorf("error querying factors: %v", err)
	}
	defer rows.Close()

	var factors []Factor
	for rows.Next() {
		var factor Factor
		if err := rows.Scan(&factor.ID, &factor.Factor, &factor.MainId); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		factors = append(factors, factor)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return factors, nil
}

// InsertFactor inserts a factor into the database
func InsertFactor(db *pgx.Conn, factor Factor) error {
	query := `INSERT INTO public.factors (id, factor, mainid) VALUES ($1, $2, $3)`
	_, err := db.Exec(context.Background(), query, factor.ID, factor.Factor, factor.MainId)
	if err != nil {
		return fmt.Errorf("error inserting factor: %v", err)
	}
	return nil
}

// InsertRecord inserts a record into the database
func InsertRecord(db *pgx.Conn, perm WritePermutation) error {
	const maxRetries = 1000
	const retryDelay = 2 * time.Second

	query := `INSERT INTO public.permutations (
								 id,
								 startArray,
								 endArray,
								 packageName,
								 permName,
								 reportedToAPI,
								 processed,
								 arrayLength,
								 numberOfPermutations)
		   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	var err error
	for i := 0; i < maxRetries; i++ {
		_, err = db.Exec(
			context.Background(),
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
			_ = fmt.Errorf("error inserting record: %v", err)
			time.Sleep(retryDelay)
			continue
		} else {
			return nil
		}
	}

	return fmt.Errorf("error inserting record after %d retries: %v", maxRetries, err)
}

// RemoveItem marks a row as processed in the database
func RemoveItem(db *pgx.Conn, id string) error {
	_, err := db.Exec(context.Background(), "DELETE FROM permutations WHERE id = $1;", id)
	if err != nil {
		return fmt.Errorf("error marking row as processed: %v", err)
	}

	return nil
}

// RemoveProcessedRows removes the processed rows from the database and compacts it
func RemoveProcessedRows(db *pgx.Conn) error {
	_, err := db.Exec(context.Background(), "DELETE FROM permutations WHERE processed = true;")
	if err != nil {
		return fmt.Errorf("error deleting processed rows: %v", err)
	}

	fmt.Println("Processed rows removed.")
	return nil
}
