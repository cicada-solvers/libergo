package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"
)

type Prime struct {
	PrimeNumber string `gorm:"column:prime_number"`
	OfPrime     string `gorm:"column:of_prime"`
}

func (Prime) TableName() string {
	return "primes"
}

func InsertPrime(db *gorm.DB, primeNumber, ofPrime string) {
	prime := Prime{
		PrimeNumber: primeNumber,
		OfPrime:     ofPrime,
	}

	result := db.Create(&prime)
	if result.Error != nil {
		fmt.Printf("error inserting prime: %v\n", result.Error)
	}
}

func GetAllPrimesByOfPrime(db *gorm.DB, ofPrime string) []Prime {
	var primes []Prime
	result := db.Model(&Prime{}).Where("of_prime = ?", ofPrime).Order("prime_number ASC").Find(&primes)
	if result.Error != nil {
		fmt.Printf("error querying primes: %v\n", result.Error)
	}
	return primes
}

func DeleteAllPrimesByOfPrime(db *gorm.DB, ofPrime string) {
	result := db.Where("of_prime = ?", ofPrime).Delete(&Prime{})
	if result.Error != nil {
		fmt.Printf("error deleting primes: %v\n", result.Error)
	}
}
