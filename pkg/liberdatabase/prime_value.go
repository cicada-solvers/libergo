package liberdatabase

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrimeValue struct {
	Id         string `gorm:"column:id"`
	PrimeValue string `gorm:"index:idx_prime,unique"`
	BitLength  int    `gorm:"index:idx_bit_length"`
}

func AddPrimeValue(db *gorm.DB, primeValue string, bitLength int) PrimeValue {
	pv := PrimeValue{
		Id:         uuid.New().String(),
		PrimeValue: primeValue,
		BitLength:  bitLength,
	}

	db.Create(&pv)
	return pv
}

func GetPrimeListLessThanValue(db *gorm.DB, value int64) []int64 {
	//var pvs []PrimeValue
	var retval []int64
	//db.Where("prime_value <= ?", value).Find(&pvs)
	//
	//for _, pv := range pvs {
	//	retval = append(retval, pv.PrimeValue)
	//}

	return sortValuesAscending(retval)
}

func sortValuesAscending(list []int64) []int64 {
	for i := 0; i < len(list)-1; i++ {
		for j := 0; j < len(list)-i-1; j++ {
			if list[j] > list[j+1] {
				list[j], list[j+1] = list[j+1], list[j]
			}
		}
	}
	return list
}
