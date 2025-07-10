package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"
)

type SentenceProb struct {
	gorm.Model
	FileName    string  `gorm:"index:idx_file_name"`
	Sentence    string  `gorm:"column:sentence"`
	Probability float64 `gorm:"column:probability"`
	GemValue    int64   `gorm:"column:gem_value"`
}

func AddSentenceProbRecord(db *gorm.DB, record SentenceProb) error {
	result := db.Create(&record)
	if result.Error != nil {
		return fmt.Errorf("error inserting sentence probability record: %v", result.Error)
	}

	return nil
}
