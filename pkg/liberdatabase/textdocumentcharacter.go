package liberdatabase

import (
	"gorm.io/gorm"
)

type TextDocumentCharacter struct {
	gorm.Model
	TextDocumentId int64  `gorm:"column:document_id;not null"`
	Character      string `gorm:"column:character;not null"`
	Count          int64  `gorm:"column:counter;not null"`
}

func (TextDocumentCharacter) TableName() string {
	return "public.text_document_characters"
}
