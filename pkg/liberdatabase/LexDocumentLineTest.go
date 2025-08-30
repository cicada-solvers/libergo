package liberdatabase

import "gorm.io/gorm"

type LexDocumentLineTest struct {
	LexDocumentLineTestId string `gorm:"column:lex_doc_line_test_id"`
	LexDocumentLineId     string `gorm:"column:lex_doc_line_id"`
	LexDocumentId         string `gorm:"column:lex_doc_id"`
	RelexifiedText        string `gorm:"column:relexified_text"`
	TextScore             int    `gorm:"column:text_score"`
}

func (LexDocumentLineTest) TableName() string {
	return "lex_document_line_tests"
}

func AddLexDocumentLineTest(db *gorm.DB, lexDocumentLineTest LexDocumentLineTest) {
	db.Create(&lexDocumentLineTest)
	return
}

func UpdateLexDocumentLineTest(db *gorm.DB, lexDocumentLineTestId string, textScore int) {
	db.Model(&LexDocumentLineTest{}).Where("lex_doc_line_test_id = ?", lexDocumentLineTestId).Update("text_score", textScore)
	return
}

func DeleteLexDocumentLineTest(db *gorm.DB, lexDocumentLineTestId string) {
	db.Where("lex_doc_line_test_id = ?", lexDocumentLineTestId).Delete(&LexDocumentLineTest{})
}

func GetAllLexDocumentLineTestsByDocumentLineId(db *gorm.DB, documentId string) []LexDocumentLineTest {
	var lexDocumentLineTests []LexDocumentLineTest
	db.Where("lex_document_line_id = ?", documentId).Find(&lexDocumentLineTests)
	return lexDocumentLineTests
}

func DeleteLexDocumentLineTextByDocumentId(db *gorm.DB, documentId string) {
	db.Where("lex_doc_id = ?", documentId).Delete(&LexDocumentLineTest{})
	return
}
