package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"
)

// TextDocument represents a text document in the database
type TextDocument struct {
	gorm.Model
	FileName string `gorm:"column:file_name;not null"`
}

// TableName sets the table name for the TextDocument struct
func (TextDocument) TableName() string {
	return "public.text_documents"
}

// InsertTextDocument inserts a new TextDocument into the database
func InsertTextDocument(db *gorm.DB, doc *TextDocument) (uint, error) {
	txn := db.Create(doc)
	if txn.Error != nil {
		fmt.Println("Error inserting text document: ", txn.Error)
		return 0, txn.Error
	}

	return doc.ID, nil
}

// GetTextDocumentByName retrieves a TextDocument by its file name
func GetTextDocumentByName(db *gorm.DB, fileName string) (*TextDocument, error) {
	var doc TextDocument
	err := db.Model(&TextDocument{}).Where("file_name = ?", fileName).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}
