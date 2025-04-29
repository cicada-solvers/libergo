package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"
)

type PrimeNumRecord struct {
	gorm.Model
	Number             int64 `gorm:"column:num"`
	IsPrime            bool  `gorm:"column:is_prime"`
	IsPrimeDuration    int64 `gorm:"column:is_prime_duration"`
	IsPtpPrime         bool  `gorm:"column:is_ptp_prime"`
	IsPtpPrimeDuration int64 `gorm:"column:is_ptp_prime_duration"`
}

func AddPrimeNumRecord(db *gorm.DB, record PrimeNumRecord) error {
	result := db.Create(&record)
	if result.Error != nil {
		return fmt.Errorf("error inserting prime number record: %v", result.Error)
	}

	return nil
}
