package liberdatabase

import "gorm.io/gorm"

type LexDocumentWord struct {
	LexDocumentWordId string  `gorm:"column:lex_doc_word_id"`
	LexDocumentId     string  `gorm:"column:lex_doc_id"`
	Word              string  `gorm:"column:word"`
	WordPosition      int64   `gorm:"column:word_position"`
	WordCount         int64   `gorm:"column:word_count"`
	WordPercentage    float64 `gorm:"column:word_percentage"`
}

func (LexDocumentWord) TableName() string {
	return "lex_document_words"
}

func AddLexDocumentWord(db *gorm.DB, words []LexDocumentWord) {
	db.Create(&words)
	return
}

func GetAllLexDocumentWordsByDocumentId(db *gorm.DB, documentId string) []LexDocumentWord {
	var words []LexDocumentWord
	db.Where("lex_document_id = ?", documentId).Find(&words)
	return words
}
