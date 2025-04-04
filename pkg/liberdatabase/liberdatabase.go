package liberdatabase

import (
	"config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

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
	dbCreateError := conn.AutoMigrate(&Permutation{})
	if dbCreateError != nil {
		fmt.Printf("Error creating Factor table: %v\n", dbCreateError)
	}
	dbCreateError = conn.AutoMigrate(&FoundHashes{})
	if dbCreateError != nil {
		fmt.Printf("Error creating Factor table: %v\n", dbCreateError)
	}
	dbCreateError = conn.AutoMigrate(&FoundHashes{})
	if dbCreateError != nil {
		fmt.Printf("Error creating table: %v\n", dbCreateError)
	}

	return conn, nil
}

func InitTables() (*gorm.DB, error) {
	// Now we need to put in our migrations.
	conn, err := InitConnection()
	if err != nil {
		return nil, err
	}

	// Migrate the schemas
	dropError := conn.Migrator().DropTable(&Permutation{})
	if dropError != nil {
		fmt.Printf("Error dropping table: %v\n", dropError)
	}
	dbCreateError := conn.AutoMigrate(&Permutation{})
	if dbCreateError != nil {
		fmt.Printf("Error creating table: %v\n", dbCreateError)
	}
	dropError = conn.Migrator().DropTable(&FoundHashes{})
	if dropError != nil {
		fmt.Printf("Error dropping table: %v\n", dropError)
	}
	dbCreateError = conn.AutoMigrate(&FoundHashes{})
	if dbCreateError != nil {
		fmt.Printf("Error creating table: %v\n", dbCreateError)
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

// InitSQLiteConnection initializes the SQLite database
func InitSQLiteConnection() (*gorm.DB, error) {
	fldrPath, err := config.GetConfigFolderPath()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}

	databasePath := filepath.Join(fldrPath, "/libergodb.db")

	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening SQLite database: %v", err)
	}
	return db, nil
}

func InitMySQLConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s%d%s", "runedonkey:dpasswd@tcp(localhost:", 3306, ")/wordsdb?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening MySQL database: %v", err)
	}
	return db, nil
}

func InitPrimesConnection() (*gorm.DB, error) {
	fldrPath, err := config.GetConfigFolderPath()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}

	databasePath := filepath.Join(fldrPath, "/primes.db")

	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening SQLite database: %v", err)
	}
	return db, nil
}

func InitPrimeTables() error {
	// Now we need to put in our migrations.
	conn, err := InitPrimesConnection()
	if err != nil {
		return nil
	}

	// Migrate the schemas
	dbCreateError := conn.AutoMigrate(&Prime{})
	if dbCreateError != nil {
		fmt.Printf("Error creating table: %v\n", dbCreateError)
	}

	return nil
}

// InitSQLiteTables initializes the SQLite database tables
func InitSQLiteTables() error {
	// Now we need to put in our migrations.
	conn, err := InitSQLiteConnection()
	if err != nil {
		return nil
	}

	// Migrate the schemas
	dbCreateError := conn.AutoMigrate(&Factor{})
	if dbCreateError != nil {
		fmt.Printf("Error creating table: %v\n", dbCreateError)
	}

	return nil
}

func InitMySqlTables() error {
	// Now we need to put in our migrations.
	conn, err := InitMySQLConnection()
	if err != nil {
		return nil
	}

	// Migrate the schemas
	dbCreateError := conn.AutoMigrate(&SentenceRecord{})
	if dbCreateError != nil {
		fmt.Printf("Error creating table: %v\n", dbCreateError)
	}

	closeError := CloseConnection(conn)
	if closeError != nil {
		return closeError
	}

	return nil
}

// CloseConnection closes the database connection
func CloseConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("error getting database instance: %v", err)
	}
	return sqlDB.Close()
}
