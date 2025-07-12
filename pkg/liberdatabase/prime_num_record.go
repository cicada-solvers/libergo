package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"
)

// PrimeNumRecord represents a database record for prime number metadata and performance metrics.
type PrimeNumRecord struct {
	gorm.Model
	Number             int64 `gorm:"column:num"`
	IsPrime            bool  `gorm:"column:is_prime"`
	IsPrimeDuration    int64 `gorm:"column:is_prime_duration"`
	IsPtpPrime         bool  `gorm:"column:is_ptp_prime"`
	IsPtpPrimeDuration int64 `gorm:"column:is_ptp_prime_duration"`
}

// AddPrimeNumRecord inserts a PrimeNumRecord into the database and returns an error if the operation fails.
func AddPrimeNumRecord(db *gorm.DB, record PrimeNumRecord) error {
	result := db.Create(&record)
	if result.Error != nil {
		return fmt.Errorf("error inserting prime number record: %v", result.Error)
	}

	return nil
}
