package liberdatabase

import (
	"conversion"
	"fmt"
	"gorm.io/gorm"
)

// Permutation represents a permutation entry in the database
type Permutation struct {
	ID                   string `gorm:"column:id"`
	StartArray           string `gorm:"column:start_array"`
	EndArray             string `gorm:"column:end_array"`
	PackageName          string `gorm:"column:package_name"`
	PermName             string `gorm:"column:perm_name"`
	ReportedToAPI        bool   `gorm:"column:reported_to_api"`
	Processed            bool   `gorm:"column:processed"`
	ArrayLength          int    `gorm:"column:array_length"`
	NumberOfPermutations int64  `gorm:"column:number_of_permutations"`
}

// ReadPermutation represents a permutation entry in the database
type ReadPermutation struct {
	ID                   string
	StartArray           []byte
	EndArray             []byte
	PackageName          string
	PermName             string
	ReportedToAPI        bool
	Processed            bool
	ArrayLength          int
	NumberOfPermutations int64
}

func (Permutation) TableName() string {
	return "public.permutations"
}

// GetByteArrayRange retrieves the unprocessed byte array ranges from the database
func GetByteArrayRange(db *gorm.DB) (*ReadPermutation, error) {
	var perm Permutation
	result := db.Model(&Permutation{}).Limit(1).Find(&perm)
	if result.Error != nil {
		return nil, fmt.Errorf("error querying row: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, nil // No more rows to process
	}

	startArray, err := conversion.ConvertToByteArray(perm.StartArray)
	if err != nil {
		return nil, fmt.Errorf("error converting start array: %v", err)
	}

	endArray, err := conversion.ConvertToByteArray(perm.EndArray)
	if err != nil {
		return nil, fmt.Errorf("error converting end array: %v", err)
	}

	return &ReadPermutation{
		ID:                   perm.ID,
		StartArray:           startArray,
		EndArray:             endArray,
		PackageName:          perm.PackageName,
		PermName:             perm.PermName,
		ReportedToAPI:        perm.ReportedToAPI,
		Processed:            perm.Processed,
		ArrayLength:          perm.ArrayLength,
		NumberOfPermutations: perm.NumberOfPermutations,
	}, nil
}

// GetByteArrayRanges retrieves the unprocessed byte array ranges from the database
func GetByteArrayRanges(db *gorm.DB) ([]ReadPermutation, error) {
	var results []ReadPermutation
	var permutations []Permutation

	result := db.Model(&Permutation{}).Where("number_of_permutations = ?", 1).Limit(25000000).Find(&permutations)
	if result.Error != nil {
		return nil, fmt.Errorf("error querying rows: %v", result.Error)
	}

	for _, perm := range permutations {
		startArray, err := conversion.ConvertToByteArray(perm.StartArray)
		if err != nil {
			return nil, fmt.Errorf("error converting start array: %v", err)
		}

		endArray, err := conversion.ConvertToByteArray(perm.EndArray)
		if err != nil {
			return nil, fmt.Errorf("error converting end array: %v", err)
		}

		results = append(results, ReadPermutation{
			ID:                   perm.ID,
			StartArray:           startArray,
			EndArray:             endArray,
			PackageName:          perm.PackageName,
			PermName:             perm.PermName,
			ReportedToAPI:        perm.ReportedToAPI,
			Processed:            perm.Processed,
			ArrayLength:          perm.ArrayLength,
			NumberOfPermutations: perm.NumberOfPermutations,
		})
	}

	return results, nil
}

// GetCountOfPermutations returns the count of rows where NumberOfPermutations = 1
func GetCountOfPermutations(db *gorm.DB) int64 {
	count := int64(0)
	result := db.Model(&Permutation{}).Where("number_of_permutations = ?", 1).Count(&count)
	if result.Error != nil {
		fmt.Printf("error counting permutations: %v\n", result.Error)
	}
	return count
}

// InsertRecord inserts a record into the database
func InsertRecord(db *gorm.DB, perm Permutation) {
	result := db.Create(&perm)
	if result.Error != nil {
		fmt.Printf("error inserting permutation: %v\n", result.Error)
	}
}

// RemoveItem marks a row as processed in the database
func RemoveItem(db *gorm.DB, id string) {
	result := db.Delete(&Permutation{}, "id = ?", id)
	if result.Error != nil {
		fmt.Printf("error deleting permutation: %v\n", result.Error)
	}
}

// RemoveProcessedRows removes the processed rows from the database and compacts it
func RemoveProcessedRows(db *gorm.DB) {
	result := db.Delete(&Permutation{}, "processed = ?", true)
	if result.Error != nil {
		fmt.Printf("error deleting permutations: %v\n", result.Error)
	}
}

func InsertBatch(db *gorm.DB, batch []Permutation) {
	result := db.Create(&batch)
	if result.Error != nil {
		fmt.Printf("Error inserting batch: %v\n", result.Error)
	}
}
