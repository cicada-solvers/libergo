package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"
)

// SentenceProb represents a record associating a sentence with its probability, file name, and calculated gem value.
type SentenceProb struct {
	gorm.Model
	FileName    string  `gorm:"index:idx_file_name"`
	Sentence    string  `gorm:"index:idx_sentence"`
	Probability float64 `gorm:"index:idx_probability"`
	GemValue    int64   `gorm:"column:gem_value"`
}

// AddSentenceProbRecord inserts a new SentenceProb record into the database and returns an error if the operation fails.
func AddSentenceProbRecord(db *gorm.DB, record SentenceProb) error {
	result := db.Create(&record)
	if result.Error != nil {
		return fmt.Errorf("error inserting sentence probability record: %v", result.Error)
	}

	return nil
}
