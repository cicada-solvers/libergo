package liberdatabase

import "gorm.io/gorm"

type WordStatistics struct {
	gorm.Model
	Word                    string  `gorm:"index:idx_word"`
	AveragePercentageOfText float64 `gorm:"column:average_percentage_of_text"`
}

func (WordStatistics) TableName() string {
	return "word_statistics"
}

func AddWordStatistics(db *gorm.DB, statistics []WordStatistics) {
	existing := db.Where("word = ?", statistics[0].Word).Find(&WordStatistics{})
	if existing.RowsAffected > 0 {
		db.Delete(&WordStatistics{}, "word = ?", statistics[0].Word)
	}

	db.Create(&statistics)
	return
}

func DeleteAllWordStatistics(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE word_statistics")
}

func GetWordByStatisticRange(db *gorm.DB, min float64, max float64) []WordStatistics {
	var statistics []WordStatistics
	db.Where("average_percentage_of_text >= ? AND average_percentage_of_text <= ?", min, max).Find(&statistics)
	return statistics
}
