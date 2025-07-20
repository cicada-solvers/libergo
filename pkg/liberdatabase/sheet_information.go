package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"
)

type SheetInformation struct {
	gorm.Model
	FileName string `gorm:"index:idx_file_name"`
	Text     string `gorm:"index:idx_file_text"`
}

// AddSheetInformation inserts a new SheetInformation record into the database and returns an error if the operation fails.
func AddSheetInformation(db *gorm.DB, record SheetInformation) error {
	result := db.Create(&record)
	if result.Error != nil {
		return fmt.Errorf("error inserting sentence record: %v", result.Error)
	}

	return nil
}
