package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"
)

// PrimeCombo represents a primecombo entry in the database
type PrimeCombo struct {
	ID        string `gorm:"column:id"`
	ValueP    string `gorm:"column:valuep"`
	ValueQ    string `gorm:"column:valueq"`
	MainId    string `gorm:"column:mainid"`
	SeqNumber int64  `gorm:"column:seqnumber"`
}

func (PrimeCombo) TableName() string {
	return "public.prime_combos"
}

// GetPrimeCombosByMainID retrieves all factors from the factors table based on the mainid.
func GetPrimeCombosByMainID(db *gorm.DB, mainId string, seqNumber int64) (*PrimeCombo, error) {
	var combo PrimeCombo

	counterResult := int64(0)
	db.Model(&PrimeCombo{}).Where("mainid = ? AND seqnumber > ?", mainId, seqNumber).Count(&counterResult)
	if counterResult == 0 {
		return nil, nil
	}

	result := db.Where("mainid = ? AND seqnumber > ?", mainId, seqNumber).First(&combo)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // No rows found, return nil
		}
		return nil, fmt.Errorf("error querying factors: %v", result.Error)
	}
	return &combo, nil
}

// InsertPrimeCombo inserts a PrimeCombo entry into the database
func InsertPrimeCombo(db *gorm.DB, combo PrimeCombo) {
	result := db.Create(&combo)
	if result.Error != nil {
		fmt.Printf("error inserting prime combo: %v\n", result.Error)
	}
}

// RemovePrimeCombosByMainID removes all prime combos from the primecombo table based on the given mainId.
func RemovePrimeCombosByMainID(db *gorm.DB, mainId string) error {
	result := db.Delete(&PrimeCombo{}, "mainid = ?", mainId)
	if result.Error != nil {
		return fmt.Errorf("error deleting prime combos: %v", result.Error)
	}
	return nil
}
