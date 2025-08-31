package liberdatabase

import "gorm.io/gorm"

type DocumentWordStatistics struct {
	gorm.Model
	FileId           string  `gorm:"index:idx_file_id"`
	Word             string  `gorm:"index:idx_word"`
	PercentageOfText float64 `gorm:"column:average_percentage_of_text"`
}

type AverageDocumentWordStatistics struct {
	Word string  `gorm:"index:word"`
	Avg  float64 `gorm:"column:avg"`
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

func GetAveraegePercentageOfTextByWord(db *gorm.DB, word string) float64 {
	var statistics AverageDocumentWordStatistics
	db.
		Table("document_word_statistics").
		Select("AVG(average_percentage_of_text)").
		Where("word = ?", word).
		Group("average_percentage_of_text").
		Group("id").First(&statistics)
	return statistics.Avg
}
