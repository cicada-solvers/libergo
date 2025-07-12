package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"
)

// Prime represents a database model for storing prime numbers and their related category.
// PrimeNumber is the string representation of the prime number.
// OfPrime is the string indicating the category or grouping of the prime.
type Prime struct {
	PrimeNumber string `gorm:"column:prime_number"`
	OfPrime     string `gorm:"column:of_prime"`
}

// TableName specifies the name of the database table associated with the Prime model.
func (Prime) TableName() string {
	return "primes"
}

// InsertPrime inserts a new prime number along with its category into the database using the provided Gorm DB instance.
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

// GetAllPrimesByOfPrime retrieves all prime records from the database filtered by the given "ofPrime" category.
// Results are ordered by the "prime_number" field in ascending order.
// Returns a slice of Prime structs containing the query results.
func GetAllPrimesByOfPrime(db *gorm.DB, ofPrime string) []Prime {
	var primes []Prime
	result := db.Model(&Prime{}).Where("of_prime = ?", ofPrime).Order("prime_number ASC").Find(&primes)
	if result.Error != nil {
		fmt.Printf("error querying primes: %v\n", result.Error)
	}
	return primes
}

// DeleteAllPrimesByOfPrime deletes all Prime records in the database where the OfPrime field matches the provided value.
func DeleteAllPrimesByOfPrime(db *gorm.DB, ofPrime string) {
	result := db.Where("of_prime = ?", ofPrime).Delete(&Prime{})
	if result.Error != nil {
		fmt.Printf("error deleting primes: %v\n", result.Error)
	}
}
