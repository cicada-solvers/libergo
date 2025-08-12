package liberdatabase

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DocumentFile represents a file in a document and its associated file type.
type DocumentFile struct {
	FileId   string `gorm:"index:idx_file_id"`
	FileName string `gorm:"index:idx_file_name"`
}

// TableName specifies the name of the database table associated with the DocumentFile model.
func (DocumentFile) TableName() string {
	return "document_files"
}

// DoesDocumentFileExist checks if a DocumentFile record exists in the database with the specified fileName.
func DoesDocumentFileExist(db *gorm.DB, fileName string) bool {
	var df DocumentFile
	result := db.Where("file_name = ?", fileName).First(&df)
	if result.Error != nil {
		return false
	}

	return true
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

// AddDocumentFile creates a new DocumentFile record in the database with the specified fileName and fileType.
// It generates a unique ID for the record and returns the newly created DocumentFile object.
func AddDocumentFile(db *gorm.DB, fileName string) DocumentFile {
	id := uuid.New().String()

	df := DocumentFile{
		FileId:   id,
		FileName: fileName,
	}

	db.Create(&df)

	return df
}
