package liberdatabase

import (
	"config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

// InitDatabase initializes the PostgreSQL database
// InitDatabase initializes the PostgreSQL database
func InitDatabase() (*gorm.DB, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}

	adminStr := cfg.AdminConnectionString

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

	// Now we need to put in our migrations.
	conn, err := InitConnection()
	if err != nil {
		return nil, err
	}

	// Migrate the schemas
	dbCreateError := conn.AutoMigrate(&Factor{})
	if dbCreateError != nil {
		fmt.Printf("Error creating Factor table: %v\n", dbCreateError)
	}
	dbCreateError = conn.AutoMigrate(&DictionaryWord{})
	if dbCreateError != nil {
		fmt.Printf("Error creating Factor table: %v\n", dbCreateError)
	}
	dbCreateError = conn.AutoMigrate(&LiberTextDocumentCharacter{})
	if dbCreateError != nil {
		fmt.Printf("Error creating Factor table: %v\n", dbCreateError)
	}
	dbCreateError = conn.AutoMigrate(&Permutation{})
	if dbCreateError != nil {
		fmt.Printf("Error creating Factor table: %v\n", dbCreateError)
	}
	dbCreateError = conn.AutoMigrate(&TextDocument{})
	if dbCreateError != nil {
		fmt.Printf("Error creating Factor table: %v\n", dbCreateError)
	}
	dbCreateError = conn.AutoMigrate(&RuneTextDocumentCharacter{})
	if dbCreateError != nil {
		fmt.Printf("Error creating Factor table: %v\n", dbCreateError)
	}
	dbCreateError = conn.AutoMigrate(&TextDocumentCharacter{})
	if dbCreateError != nil {
		fmt.Printf("Error creating Factor table: %v\n", dbCreateError)
	}
	dbCreateError = conn.AutoMigrate(&FoundHashes{})
	if dbCreateError != nil {
		fmt.Printf("Error creating Factor table: %v\n", dbCreateError)
	}

	return conn, nil
}

// InitConnection initializes the PostgreSQL database
func InitConnection() (*gorm.DB, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}

	connStr := cfg.GeneralConnectionString

	dsn := connStr
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
