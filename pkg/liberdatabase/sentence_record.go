package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"
)

type SentenceRecord struct {
	gorm.Model
	FileName     string `gorm:"column:file_name"`
	DictSentence string `gorm:"column:dict_sentence"`
}

func AddSentenceRecord(db *gorm.DB, record SentenceRecord) error {
	result := db.Create(&record)
	if result.Error != nil {
		return fmt.Errorf("error inserting sentence record: %v", result.Error)
	}

	return nil
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
