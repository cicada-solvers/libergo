package liberdatabase

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AdvancedNumberFactors represents a struct for storing detailed factors of a number and their calculated percentages.
// Id is the unique identifier of the record.
// AdvancedNumberInformationId links to a related AdvancedNumberInformation.
// Factor is the numerical factor of the associated number.
// FactorPosition is the order or position of the factor in the factorization.
// PercentFromSquareRoot is the percentage of the factor relative to the square root of the number.
// PercentFromNumber is the percentage of the factor relative to the number itself.
// PercentFromTwo is the percentage of the factor relative to the number two.
// PercentFromMiddle is the percentage of the factor relative to a midpoint calculation.
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

// AddAdvancedNumberFactors adds a record of advanced number factors with associated percentages to the database.
// It links the factor to a specific advanced number information record by its ID and calculates various percentage metrics.
// Parameters include database instance, advanced number information ID, factor details, and percentage calculations.
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
