package liberdatabase

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrimeValue struct {
	Id         string `gorm:"column:id"`
	PrimeValue int64  `gorm:"column:prime_value"`
}

func AddPrimeValue(db *gorm.DB, primeValue int64) PrimeValue {
	pv := PrimeValue{
		Id:         uuid.New().String(),
		PrimeValue: primeValue,
	}

	db.Create(&pv)
	return pv
}
