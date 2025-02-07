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
