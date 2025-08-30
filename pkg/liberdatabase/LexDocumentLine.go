package liberdatabase

import "gorm.io/gorm"

type LexDocumentLine struct {
	LexDocumentLineId string `gorm:"column:lex_doc_line_id"`
	LexDocumentId     string `gorm:"column:lex_doc_id"`
	Line              string `gorm:"column:line"`
	LinePosition      int64  `gorm:"column:line_position"`
}

func (LexDocumentLine) TableName() string {
	return "lex_document_lines"
}

func AddLexDocumentLine(db *gorm.DB, lexDocumentLine LexDocumentLine) {
	db.Create(&lexDocumentLine)
	return
}

func GetAllLexDocumentLinesByDocumentId(db *gorm.DB, documentId string) []LexDocumentLine {
	var lexDocumentLines []LexDocumentLine
	db.Where("lex_document_id = ?", documentId).Find(&lexDocumentLines)
	return lexDocumentLines
}

func DeleteLexDocumentLinesByDocumentId(db *gorm.DB, documentId string) {
	db.Where("lex_doc_id = ?", documentId).Delete(&LexDocumentLine{})
	return
}
