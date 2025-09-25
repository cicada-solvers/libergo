package liberdatabase

import "gorm.io/gorm"

// GoldbachAddend represents the addends of a Goldbach pair for an even number in the database.
// GoldbachId is the unique identifier linking this addend pair to a specific even number.
// AddendOne is the first prime addend of the Goldbach pair.
// AddendTwo is the second prime addend of the Goldbach pair.
// AddendThree is the third prime addend of the Goldbach pair.
type GoldbachAddend struct {
	GoldbachNumber int64 `gorm:"column:goldbach_number"`
	AddendOne      int64 `gorm:"column:addend_one"`
	AddendTwo      int64 `gorm:"column:addend_two"`
	AddendThree    int64 `gorm:"column:addend_three"`
	SetNumber      int   `gorm:"column:set_number"`
}

// AddGoldbachAddends inserts a list of Goldbach addends into the database and returns the inserted addends.
func AddGoldbachAddends(db *gorm.DB, addends []GoldbachAddend) []GoldbachAddend {
	db.Create(&addends)

	return addends
}
