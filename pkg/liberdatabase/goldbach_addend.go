package liberdatabase

import "gorm.io/gorm"

type GoldbachAddend struct {
	GoldbachId string `gorm:"column:goldbach_id"`
	AddendOne  int64  `gorm:"column:addend_one"`
	AddendTwo  int64  `gorm:"column:addend_two"`
}

func AddGoldbachAddend(db *gorm.DB, goldbachId string, addendOne int64, addendTwo int64) GoldbachAddend {
	goldbachAddend := GoldbachAddend{
		GoldbachId: goldbachId,
		AddendOne:  addendOne,
		AddendTwo:  addendTwo,
	}

	db.Create(&goldbachAddend)

	return goldbachAddend
}
