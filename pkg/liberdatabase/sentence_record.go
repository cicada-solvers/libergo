package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"
)

// SentenceRecord represents a database record containing a sentence, its associated file name, and additional metadata.
type SentenceRecord struct {
	gorm.Model
	FileName     string `gorm:"index:idx_file_name"`
	DictSentence string `gorm:"index:idx_sentence"`
	GemValue     int64  `gorm:"column:gem_value"`
	IsPrime      bool   `gorm:"column:is_prime"`
}

// AddSentenceRecord inserts a slice of SentenceRecord structures into the database using Gorm's Create method.
// Returns an error if the insertion process fails.
func AddSentenceRecord(db *gorm.DB, records []SentenceRecord) error {
	result := db.Create(&records)
	if result.Error != nil {
		return fmt.Errorf("error inserting sentence record: %v", result.Error)
	}

	return nil
}

// GetRecordCountByFileName returns the count of SentenceRecord entries in the database matching the specified fileName.
func GetRecordCountByFileName(db *gorm.DB, fileName string) (int64, error) {
	var count int64
	result := db.Model(&SentenceRecord{}).Where("file_name = ?", fileName).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("error retrieving record count: %v", result.Error)
	}
	return count, nil
}

// GetAllFileNames retrieves all distinct file names from the database and maps each file name to its associated record count.
func GetAllFileNames(db *gorm.DB) (map[string]int, error) {
	var fileNames []string
	var fileCountMap = make(map[string]int)
	result := db.Model(&SentenceRecord{}).Select("DISTINCT file_name").Find(&fileNames)
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving file names: %v", result.Error)
	}

	for _, fileName := range fileNames {
		count, err := GetRecordCountByFileName(db, fileName)
		if err != nil {
			return nil, fmt.Errorf("error retrieving record count: %v", err)
		}
		fileCountMap[fileName] = int(count)
	}

	return fileCountMap, nil
}

// GetTopMillionSentenceRecords retrieves up to one million SentenceRecord entries for a specified file name from the database.
// Returns the records and any encountered error.
func GetTopMillionSentenceRecords(db *gorm.DB, fileName string) ([]SentenceRecord, error) {
	var records []SentenceRecord
	result := db.Model(&SentenceRecord{}).Where("file_name = ?", fileName).Limit(1000000).Find(&records)
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving top million sentence records: %v", result.Error)
	}
	return records, nil
}

// RemoveMillionSentenceRecords deletes a batch of SentenceRecord entries from the database and returns an error if any occur.
func RemoveMillionSentenceRecords(db *gorm.DB, records []SentenceRecord) error {
	for _, record := range records {
		result := db.Delete(&record)
		if result.Error != nil {
			return fmt.Errorf("error removing sentence record: %v", result.Error)
		}
	}
	return nil
}

// RemoveAllSentenceRecords removes all records of type SentenceRecord from the database using the provided gorm.DB instance.
// Returns an error if the deletion operation fails.
func RemoveAllSentenceRecords(db *gorm.DB) error {
	result := db.Delete(&SentenceRecord{})
	if result.Error != nil {
		return fmt.Errorf("error removing all sentence records: %v", result.Error)
	}
	return nil
}
