package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"
)

type PrimeNumRecord struct {
	gorm.Model
	Number                 int64  `gorm:"column:num"`
	IsPrime                bool   `gorm:"column:is_prime"`
	NumberCountBeforePrime int64  `gorm:"column:number_count_before_prime"`
	PrimeFactorCount       int64  `gorm:"column:prime_factor_count"`
	PrimeFactors           string `gorm:"column:prime_factors"`
}

func AddPrimeNumRecord(db *gorm.DB, record PrimeNumRecord) error {
	result := db.Create(&record)
	if result.Error != nil {
		return fmt.Errorf("error inserting prime number record: %v", result.Error)
	}

	return nil
}
