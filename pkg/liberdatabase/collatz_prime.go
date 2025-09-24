package liberdatabase

import "gorm.io/gorm"

type CollatzPrime struct {
	gorm.Model
	Number                    int64  `gorm:"index:idx_number"`
	IsPrime                   bool   `gorm:"column:is_prime"`
	CollatzPrimeLength        int64  `gorm:"column:collatz_prime_length"`
	CollatzSequence           string `gorm:"column:collatz_sequence"`
	CollatzPrimeLengthIsPrime bool   `gorm:"column:collatz_prime_length_is_prime"`
}

func (CollatzPrime) TableName() string {
	return "collatz_primes"
}

func AddCollatzPrimes(db *gorm.DB, collatzPrimes []CollatzPrime) {
	db.Create(&collatzPrimes)
	return
}
