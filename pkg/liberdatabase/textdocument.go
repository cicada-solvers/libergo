package liberdatabase

import "gorm.io/gorm"

type TextDocument struct {
	gorm.Model
	FileName string `gorm:"column:file_name;not null"`
}

func (TextDocument) TableName() string {
	return "public.text_documents"
}
