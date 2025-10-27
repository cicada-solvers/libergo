package liberdatabase

import "gorm.io/gorm"

type DocumentCharacter struct {
	FileId           string  `gorm:"index:idx_character"`
	Character        string  `gorm:"index:idx_character"`
	CharacterCount   int64   `gorm:"column:character_count"`
	CharacterType    string  `gorm:"index:idx_character"`
	PercentageOfText float64 `gorm:"column:percentage_of_text"`
}

func (DocumentCharacter) TableName() string {
	return "document_characters"
}

func GetDocumentCharactersByFileIdAndCharacterType(db *gorm.DB, fileId string, characterType string) []DocumentCharacter {
	var characters = make([]DocumentCharacter, 0)
	db.Where("file_id = ? AND character_type = ?", fileId, characterType).Find(&characters)
	return characters
}

func AddDocumentCharacters(db *gorm.DB, characters []DocumentCharacter) {
	db.Create(&characters)
	return
}

func UpdatePercentageOfText(db *gorm.DB, fileId string, character string, percentageOfText float64) {
	db.Model(&DocumentCharacter{}).Where("file_id = ? AND character = ?", fileId, character).Update("percentage_of_text", percentageOfText)
}

func DeleteCharactersByFileId(db *gorm.DB, fileId string) {
	db.Where("file_id = ?", fileId).Delete(&DocumentCharacter{})
}
