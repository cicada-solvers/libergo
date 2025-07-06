package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"
)

type SentenceRecord struct {
	gorm.Model
	FileName     string `gorm:"column:file_name"`
	DictSentence string `gorm:"column:dict_sentence"`
	GemValue     int64  `gorm:"column:gem_value"`
	IsPrime      bool   `gorm:"column:is_prime"`
}

func AddSentenceRecord(db *gorm.DB, records []SentenceRecord) error {
	result := db.Create(&records)
	if result.Error != nil {
		return fmt.Errorf("error inserting sentence record: %v", result.Error)
	}

	return nil
}

func GetRecordCountByFileName(db *gorm.DB, fileName string) (int64, error) {
	var count int64
	result := db.Model(&SentenceRecord{}).Where("file_name = ?", fileName).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("error retrieving record count: %v", result.Error)
	}
	return count, nil
}

func GetAllFileNames(db *gorm.DB) ([]string, error) {
	var fileNames []string
	result := db.Model(&SentenceRecord{}).Select("DISTINCT file_name").Find(&fileNames)
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving file names: %v", result.Error)
	}
	return fileNames, nil
}

func GetTopMillionSentenceRecords(db *gorm.DB, fileName string) ([]SentenceRecord, error) {
	var records []SentenceRecord
	result := db.Model(&SentenceRecord{}).Where("file_name = ?", fileName).Limit(1000000).Find(&records)
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving top million sentence records: %v", result.Error)
	}
	return records, nil
}

func RemoveMillionSentenceRecords(db *gorm.DB, records []SentenceRecord) error {
	for _, record := range records {
		result := db.Delete(&record)
		if result.Error != nil {
			return fmt.Errorf("error removing sentence record: %v", result.Error)
		}
	}
	return nil
}

func RemoveAllSentenceRecords(db *gorm.DB) error {
	result := db.Delete(&SentenceRecord{})
	if result.Error != nil {
		return fmt.Errorf("error removing all sentence records: %v", result.Error)
	}
	return nil
}
