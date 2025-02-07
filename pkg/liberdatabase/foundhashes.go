package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"
)

type FoundHashes struct {
	gorm.Model
	IntegerArray     string `gorm:"column:integer_array"`
	HashingAlgorithm string `gorm:"column:hashing_algorithm"`
}

func (FoundHashes) TableName() string {
	return "public.found_hashes"
}

// InsertFoundHash inserts a new FoundHashes record into the database
func InsertFoundHash(integerArray string, hashingAlgorithm string) error {
	db, _ := InitConnection()

	foundHash := FoundHashes{
		IntegerArray:     integerArray,
		HashingAlgorithm: hashingAlgorithm,
	}

	result := db.Create(&foundHash)
	if result.Error != nil {
		return fmt.Errorf("error inserting found hash: %v", result.Error)
	}

	return nil
}

// GetAllFoundHashes retrieves all FoundHashes records from the database and prints them to the console
func GetAllFoundHashes() error {
	db, err := InitConnection()
	if err != nil {
		return fmt.Errorf("error initializing database connection: %v", err)
	}

	var foundHashes []FoundHashes
	result := db.Find(&foundHashes)
	if result.Error != nil {
		return fmt.Errorf("error retrieving found hashes: %v", result.Error)
	}

	fmt.Println("Found Hashes:")
	if len(foundHashes) == 0 {
		fmt.Println("No found hashes found.")
		return nil
	}

	for _, foundHash := range foundHashes {
		fmt.Printf("ID: %d, IntegerArray: %s, HashingAlgorithm: %s, CreatedAt: %s, UpdatedAt: %s\n",
			foundHash.ID, foundHash.IntegerArray, foundHash.HashingAlgorithm, foundHash.CreatedAt, foundHash.UpdatedAt)
	}

	return nil
}
