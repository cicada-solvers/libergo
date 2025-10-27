package liberdatabase

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DocumentFile represents a file in a document and its associated file type.
type DocumentFile struct {
	FileId              string `gorm:"index:idx_file_id"`
	FileName            string `gorm:"index:idx_file_name"`
	WordCount           int64  `gorm:"column:word_count"`
	TotalCharacterCount int64  `gorm:"column:total_character_count"`
	TotalRuneCount      int64  `gorm:"column:total_rune_count"`
}

// TableName specifies the name of the database table associated with the DocumentFile model.
func (DocumentFile) TableName() string {
	return "document_files"
}

// DoesDocumentFileExist checks if a DocumentFile record exists in the database with the specified fileName.
// Uses an optimized query that only checks for existence without retrieving the full record.
func DoesDocumentFileExist(db *gorm.DB, fileName string) bool {
	var exists bool
	result := db.Model(&DocumentFile{}).
		Select("count(*) > 0").
		Where("file_name = ?", fileName).
		Find(&exists)

	if result.Error != nil {
		return false
	}

	return exists
}

// GetDocumentFile retrieves a DocumentFile record by its fileName from the database. Returns the record or an error if not found.
func GetDocumentFile(db *gorm.DB, fileName string) (DocumentFile, error) {
	var df DocumentFile
	result := db.Where("file_name = ?", fileName).First(&df)
	if result.Error != nil {
		return df, result.Error
	}
	return df, nil
}

func GetDocumentFileById(db *gorm.DB, fileId string) (DocumentFile, error) {
	var df DocumentFile
	result := db.Where("file_id = ?", fileId).First(&df)
	if result.Error != nil {
		return df, result.Error
	}

	return df, nil
}

// GetAllDocumentFiles retrieves all DocumentFile records from the database.
func GetAllDocumentFiles(db *gorm.DB) []DocumentFile {
	var df []DocumentFile
	db.Find(&df)
	return df
}

// AddDocumentFile creates a new DocumentFile record in the database with the specified fileName and fileType.
// It generates a unique ID for the record and returns the newly created DocumentFile object.
func AddDocumentFile(db *gorm.DB, fileName string) DocumentFile {
	id := uuid.New().String()

	df := DocumentFile{
		FileId:    id,
		FileName:  fileName,
		WordCount: 0,
	}

	db.Create(&df)

	return df
}

func UpdateDocumentWordCount(db *gorm.DB, fileId string, wordCount int64) {
	db.Model(&DocumentFile{}).Where("file_id = ?", fileId).Update("word_count", wordCount)
	return
}

func UpdateTotalCharacterCount(db *gorm.DB, fileId string, totalCharacterCount int64) {
	db.Model(&DocumentFile{}).Where("file_id = ?", fileId).Update("total_character_count", totalCharacterCount)
}

func UpdateTotalRuneCount(db *gorm.DB, fileId string, totalRuneCount int64) {
	db.Model(&DocumentFile{}).Where("file_id = ?", fileId).Update("total_rune_count", totalRuneCount)
}
