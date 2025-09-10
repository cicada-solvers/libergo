package liberdatabase

import "gorm.io/gorm"

type OutputData struct {
	gorm.Model
	DocId string `gorm:"index:idx_output_data_doc_id"`
	Data  string `gorm:"index:idx_output_data_doc_id"`
	Score int64  `gorm:"index:idx_output_data_doc_id"`
}

func (OutputData) TableName() string {
	return "output_data"
}

func AddOutputData(db *gorm.DB, outputData OutputData) {
	db.Create(&outputData)
	return
}

func GetOutputDataByDocId(db *gorm.DB, docId string) []OutputData {
	var outputData []OutputData
	db.Where("doc_id = ?", docId).Find(&outputData)
	return outputData
}
