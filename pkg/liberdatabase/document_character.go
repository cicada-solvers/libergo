package liberdatabase

import "gorm.io/gorm"

type DocumentCharacter struct {
	FileId           string `gorm:"index:idx_file_id"`
	Character        string `gorm:"index:idx_character"`
	CharacterCount   int64  `gorm:"column:character_count"`
	CharacterType    string `gorm:"column:character_type"`
	PercentageOfText float64
}

func (DocumentCharacter) TableName() string {
	return "document_characters"
}

func AddDocumentCharacters(db *gorm.DB, characters []DocumentCharacter) {
	db.Create(&characters)
	return
}

func DeleteCharactersByFileId(db *gorm.DB, fileId string) {
	db.Where("file_id = ?", fileId).Delete(&DocumentCharacter{})
}
