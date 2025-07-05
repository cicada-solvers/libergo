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
