package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"
)

type PrimeNumRecord struct {
	gorm.Model
	Number                 string `gorm:"column:num"`
	NumberCountBeforePrime string `gorm:"column:number_count_before_prime"`
	NumberIsPrime          bool   `gorm:"column:is_prime"`
	NumberFactorSize       int64  `gorm:"column:factorsize"`
}

func AddPrimeNumRecord(db *gorm.DB, record PrimeNumRecord) error {
	result := db.Create(&record)
	if result.Error != nil {
		return fmt.Errorf("error inserting prime number record: %v", result.Error)
	}

	return nil
}
