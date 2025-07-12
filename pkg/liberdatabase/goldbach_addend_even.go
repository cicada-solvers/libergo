package liberdatabase

import "gorm.io/gorm"

// GoldbachAddendEven represents the addends of a Goldbach pair for an even number in the database.
// GoldbachId is the unique identifier linking this addend pair to a specific even number.
// AddendOne is the first prime addend of the Goldbach pair.
// AddendTwo is the second prime addend of the Goldbach pair.
type GoldbachAddendEven struct {
	GoldbachId string `gorm:"column:goldbach_id"`
	AddendOne  int64  `gorm:"column:addend_one"`
	AddendTwo  int64  `gorm:"column:addend_two"`
}

// AddGoldbachAddends inserts a list of Goldbach addends into the database and returns the inserted addends.
func AddGoldbachAddends(db *gorm.DB, addends []GoldbachAddendEven) []GoldbachAddendEven {
	db.Create(&addends)

	return addends
}
