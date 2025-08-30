package liberdatabase

import (
	"gorm.io/gorm"
)

// DocumentWord represents a word in a document and its associated word count.
type DocumentWord struct {
	gorm.Model
	Word      string `gorm:"index:idx_word"`
	FileId    string `gorm:"index:idx_file_id"`
	WordCount int64  `gorm:"column:word_count"`
}

// TableName specifies the database table name for the DocumentWord model.
func (DocumentWord) TableName() string {
	return "document_words"
}

// AddDocumentWord inserts a new DocumentWord record into the database.
func AddDocumentWord(db *gorm.DB, words []DocumentWord) {
	db.Create(&words)
}

func GetDistinctWords(db *gorm.DB, fileId string) []DocumentWord {
	var words []DocumentWord
	db.Where("file_id = ?", fileId).Find(&words)
	return words
}

func DeleteWordsByFileId(db *gorm.DB, fileId string) {
	db.Where("file_id = ?", fileId).Delete(&DocumentWord{})
	return
}

func GetAllDistinctWords(db *gorm.DB, minId uint) []DocumentWord {
	var words []DocumentWord
	db.Where("id > ?", minId).Order("id ASC").Find(&words).Limit(500)
	return words
}
