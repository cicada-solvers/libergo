package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"
)

type PrimeNumRecord struct {
	gorm.Model
	Number                 int64 `gorm:"column:num"`
	NumberCountBeforePrime int   `gorm:"column:number_count_before_prime"`
}

func AddPrimeNumRecord(db *gorm.DB, record PrimeNumRecord) error {
	result := db.Create(&record)
	if result.Error != nil {
		return fmt.Errorf("error inserting prime number record: %v", result.Error)
	}

	return nil
}
