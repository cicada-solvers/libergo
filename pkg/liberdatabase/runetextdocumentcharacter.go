package liberdatabase

import "gorm.io/gorm"

type RuneTextDocumentCharacter struct {
	gorm.Model
	TextDocumentId int64  `gorm:"column:document_id;not null"`
	Character      string `gorm:"column:character;not null"`
	Count          int64  `gorm:"column:counter;not null"`
}

func (RuneTextDocumentCharacter) TableName() string {
	return "public.rune_text_document_characters"
}

func InsertRuneTextDocumentCharacter(db *gorm.DB, tdc *RuneTextDocumentCharacter) (uint, error) {
	if err := db.Create(tdc).Error; err != nil {
		return 0, err
	}
	return tdc.ID, nil
}
