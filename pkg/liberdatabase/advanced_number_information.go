package liberdatabase

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdvancedNumberInformation struct {
	Id         string `gorm:"column:id"`
	Number     string `gorm:"column:number"`
	SquareRoot string `gorm:"column:square_root"`
}

func AddAdvancedNumberInformation(db *gorm.DB, number string, squareRoot string) AdvancedNumberInformation {
	advancedNumberInformation := AdvancedNumberInformation{
		Id:         uuid.New().String(),
		Number:     number,
		SquareRoot: squareRoot,
	}

	db.Create(&advancedNumberInformation)

	return advancedNumberInformation
}
