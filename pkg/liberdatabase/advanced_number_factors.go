package liberdatabase

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdvancedNumberFactors struct {
	Id                          string  `gorm:"column:id"`
	AdvancedNumberInformationId string  `gorm:"column:advanced_number_information_id"`
	Factor                      int64   `gorm:"column:factor"`
	FactorPosition              int     `gorm:"column:factor_position"`
	PercentFromSquareRoot       float64 `gorm:"column:percent_from_square_root"`
	PercentFromNumber           float64 `gorm:"column:percent_from_number"`
	PercentFromTwo              float64 `gorm:"column:percent_from_two"`
	PercentFromMiddle           float64 `gorm:"column:percent_from_middle"`
}

func AddAdvancedNumberFactors(db *gorm.DB, advancedNumberInformationId string, factor int64, factorPosition int,
	percentFromSquareRoot float64, percentFromNumber float64, percentFromTwo float64, percentFromMiddle float64) {
	advancedNumberFactors := AdvancedNumberFactors{
		Id:                          uuid.New().String(),
		AdvancedNumberInformationId: advancedNumberInformationId,
		Factor:                      factor,
		FactorPosition:              factorPosition,
		PercentFromSquareRoot:       percentFromSquareRoot,
		PercentFromNumber:           percentFromNumber,
		PercentFromTwo:              percentFromTwo,
		PercentFromMiddle:           percentFromMiddle,
	}

	db.Create(&advancedNumberFactors)
	return
}
