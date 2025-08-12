package liberdatabase

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DocumentWord represents a word in a document and its associated word count.
type DocumentWord struct {
	Id        string `gorm:"column:id"`
	Word      string `gorm:"index:idx_word"`
	FileId    string `gorm:"index:idx_file_id"`
	WordCount int64  `gorm:"column:word_count"`
}

// TableName specifies the database table name for the DocumentWord model.
func (DocumentWord) TableName() string {
	return "document_words"
}

// AddDocumentWord inserts a new DocumentWord record into the database.
func AddDocumentWord(db *gorm.DB, word, fileId string, wordCount int64) {
	id := uuid.New().String()

	dw := DocumentWord{
		Id:        id,
		Word:      word,
		FileId:    fileId,
		WordCount: wordCount,
	}

	db.Create(&dw)
}

// IncrementWordCount increments the word count for a DocumentWord record with the specified word and fileId.
func IncrementWordCount(db *gorm.DB, word, fileId string) {
	db.Model(&DocumentWord{}).Where("word = ? AND file_id = ?", word, fileId).Update("word_count", gorm.Expr("word_count + 1"))
}

// DoesWordExist checks if a DocumentWord record exists in the database with the specified word and fileId.
func DoesWordExist(db *gorm.DB, word, fileId string) bool {
	var count int64
	result := db.Model(&DocumentWord{}).Where("word = ? AND file_id = ?", word, fileId).Count(&count)
	if result.Error != nil {
		fmt.Printf("error querying words: %v\n", result.Error)
	}
	return count > 0
}
