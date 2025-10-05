package liberdatabase

import "gorm.io/gorm"

type DocumentWordStatistics struct {
	gorm.Model
	FileId           string  `gorm:"index:idx_file_id"`
	Word             string  `gorm:"index:idx_word"`
	PercentageOfText float64 `gorm:"column:percentage_of_text"`
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
	var statistics = make([]DocumentWordStatistics, 0)
	db.Where("word = ?", word).Find(&statistics)
	var wordValues = make([]float64, 0)

	for _, statistic := range statistics {
		wordValues = append(wordValues, statistic.PercentageOfText)
	}

	return calculateStatistics(wordValues)
}

func calculateStatistics(words []float64) float64 {
	total := 0.0
	for _, word := range words {
		total += word
	}

	return total / float64(len(words))
}
