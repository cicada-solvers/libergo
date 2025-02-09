package liberdatabase

import "gorm.io/gorm"

// RuneglishTextDocumentCharacter represents a character in a text document
type RuneglishTextDocumentCharacter struct {
	gorm.Model
	TextDocumentId int64  `gorm:"column:document_id;not null"`
	Character      string `gorm:"column:character;not null"`
	Count          int64  `gorm:"column:count;not null"`
}

func (RuneglishTextDocumentCharacter) TableName() string {
	return "public.runeglish_text_document_characters"
}

func InsertRuneglishTextDocumentCharacter(db *gorm.DB, tdc *RuneglishTextDocumentCharacter) (uint, error) {
	if err := db.Create(tdc).Error; err != nil {
		return 0, err
	}
	return tdc.ID, nil
}
