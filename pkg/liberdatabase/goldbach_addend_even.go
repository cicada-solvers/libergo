package liberdatabase

import "gorm.io/gorm"

type GoldbachAddendEven struct {
	GoldbachId string `gorm:"column:goldbach_id"`
	AddendOne  int64  `gorm:"column:addend_one"`
	AddendTwo  int64  `gorm:"column:addend_two"`
}

func AddGoldbachAddends(db *gorm.DB, addends []GoldbachAddendEven) []GoldbachAddendEven {
	db.Create(&addends)

	return addends
}
