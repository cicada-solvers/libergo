package liberdatabase

import "gorm.io/gorm"

type LexDocument struct {
	LexDocumentId string `gorm:"column:lex_doc_id"`
	FileName      string `gorm:"column:file_name"`
}

func (LexDocument) TableName() string {
	return "lex_documents"
}

func DoesLexDocumentExist(db *gorm.DB, fileName string) bool {
	var lexDocument LexDocument
	db.Where("file_name = ?", fileName).First(&lexDocument)
	return lexDocument.LexDocumentId != ""
}

func AddLexDocument(db *gorm.DB, lexDocument LexDocument) {
	db.Create(&lexDocument)
}

func GetAllLexDocuments(db *gorm.DB) []LexDocument {
	var lexDocuments []LexDocument
	db.Find(&lexDocuments)
	return lexDocuments
}

func GetLexDocumentByFileName(db *gorm.DB, fileName string) LexDocument {
	var lexDocument LexDocument
	db.Where("file_name = ?", fileName).First(&lexDocument)
	return lexDocument
}
