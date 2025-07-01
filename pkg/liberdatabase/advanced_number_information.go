package liberdatabase

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdvancedNumberInformation struct {
	Id         string `gorm:"column:id"`
	Number     int64  `gorm:"column:number"`
	SquareRoot int64  `gorm:"column:square_root"`
}

func AddAdvancedNumberInformation(db *gorm.DB, number, squareRoot int64) AdvancedNumberInformation {
	advancedNumberInformation := AdvancedNumberInformation{
		Id:         uuid.New().String(),
		Number:     number,
		SquareRoot: squareRoot,
	}

	db.Create(&advancedNumberInformation)

	return advancedNumberInformation
}
