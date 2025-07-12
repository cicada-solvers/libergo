package liberdatabase

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GoldbachNumberEven represents an even number associated with the Goldbach conjecture, along with its database details.
// Id is the unique identifier held as a string.
// Number stores the associated even number as int64.
// IsEven is a boolean indicating whether the number is even (always true for this type).
type GoldbachNumberEven struct {
	Id     string `gorm:"column:id"`
	Number int64  `gorm:"column:number"`
	IsEven bool   `gorm:"column:is_even"`
}

// AddGoldbachNumber adds a GoldbachNumberEven entry to the database with the given number and even status.
// Takes a gorm.DB instance, an int64 number, and a boolean indicating if the number is even.
// Returns the created GoldbachNumberEven object.
func AddGoldbachNumber(db *gorm.DB, number int64, isEven bool) GoldbachNumberEven {
	goldbachNumber := GoldbachNumberEven{
		Id:     uuid.New().String(),
		Number: number,
		IsEven: isEven,
	}

	db.Create(&goldbachNumber)

	return goldbachNumber
}
