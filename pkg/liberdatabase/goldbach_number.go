package liberdatabase

import (
	"gorm.io/gorm"
)

// GoldbachNumber represents an even number associated with the Goldbach conjecture, along with its database details.
// Id is the unique identifier held as a string.
// Number stores the associated even number as int64.
// IsEven is a boolean indicating whether the number is even (always true for this type).
type GoldbachNumber struct {
	gorm.Model
	Number  int64 `gorm:"column:number"`
	IsEven  bool  `gorm:"column:is_even"`
	IsPrime bool  `gorm:"column:is_prime"`
}

// AddGoldbachNumber adds a GoldbachNumberEven entry to the database with the given number and even status.
// Takes a gorm.DB instance, an int64 number, and a boolean indicating if the number is even.
// Returns the created GoldbachNumberEven object.
func AddGoldbachNumber(db *gorm.DB, number int64, isEven bool, isPrime bool) GoldbachNumber {
	goldbachNumber := GoldbachNumber{
		Number:  number,
		IsEven:  isEven,
		IsPrime: isPrime,
	}

	db.Create(&goldbachNumber)

	return goldbachNumber
}
