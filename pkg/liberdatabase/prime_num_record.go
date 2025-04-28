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
	ModTwoTen              int64  `gorm:"column:mod_two_ten"`
	ModTwoTenIsPrime       bool   `gorm:"column:mod_two_ten_is_prime"`
	ModTwoTenFactors       string `gorm:"column:mod_two_ten_factors"`
	ModFortyEight          int64  `gorm:"column:mod_forty_eight"`
	ModFortyEightIsPrime   bool   `gorm:"column:mod_forty_eight_is_prime"`
	ModFortyEightFactors   string `gorm:"column:mod_forty_eight_factors"`
	ModTen                 int64  `gorm:"column:mod_ten"`
	ModTenIsPrime          bool   `gorm:"column:mod_ten_is_prime"`
	ModTenFactors          string `gorm:"column:mod_ten_factors"`
}

func AddPrimeNumRecord(db *gorm.DB, record PrimeNumRecord) error {
	result := db.Create(&record)
	if result.Error != nil {
		return fmt.Errorf("error inserting prime number record: %v", result.Error)
	}

	return nil
}
