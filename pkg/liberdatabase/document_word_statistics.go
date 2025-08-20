package liberdatabase

import "gorm.io/gorm"

type DocumentWordStatistics struct {
	Id               string  `gorm:"column:id"`
	FileId           string  `gorm:"index:idx_file_id"`
	Word             string  `gorm:"index:idx_word"`
	PercentageOfText float64 `gorm:"column:average_percentage_of_text"`
}

func (DocumentWordStatistics) TableName() string {
	return "document_word_statistics"
}

func AddDocumentWordStatistics(db *gorm.DB, statistics []DocumentWordStatistics) {
	db.Create(&statistics)
}

func DeleteStatisticsByFileId(db *gorm.DB, fileId string) {
	db.Where("file_id = ?", fileId).Delete(&DocumentWordStatistics{})
	return
}
