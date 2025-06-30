package liberdatabase

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GoldbachNumber struct {
	Id     string `gorm:"column:id"`
	Number int64  `gorm:"column:number"`
	IsEven bool   `gorm:"column:is_even"`
}

func AddGoldbachNumber(db *gorm.DB, number int64, isEven bool) GoldbachNumber {
	goldbachNumber := GoldbachNumber{
		Id:     uuid.New().String(),
		Number: number,
		IsEven: isEven,
	}

	db.Create(&goldbachNumber)

	return goldbachNumber
}
