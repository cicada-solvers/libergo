package main

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"liberdatabase"
)

type AdvancedNumberInformation struct {
	Id         string `gorm:"column:id"`
	Number     string `gorm:"column:number"`
	SquareRoot string `gorm:"column:square_root"`
}

type AdvancedNumberFactors struct {
	Id                          string  `gorm:"column:id"`
	AdvancedNumberInformationId string  `gorm:"column:advanced_number_information_id"`
	Factor                      string  `gorm:"column:factor"`
	FactorPosition              int     `gorm:"column:factor_position"`
	PercentFromSquareRoot       float64 `gorm:"column:percent_from_square_root"`
	PercentFromNumber           float64 `gorm:"column:percent_from_number"`
	PercentFromTwo              float64 `gorm:"column:percent_from_two"`
	PercentFromMiddle           float64 `gorm:"column:percent_from_middle"`
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

func AddAdvancedNumberFactors(db *gorm.DB, advancedNumberInformationId string, factor string, factorPosition int,
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

func CreateDatabase() {
	db, connError := liberdatabase.InitDatabase()
	if connError != nil {
		fmt.Printf("Error initializing database connection: %v\n", connError)
		return
	}

	// Remove the old table if it exists
	dropError := db.Migrator().DropTable(&AdvancedNumberInformation{})
	if dropError != nil {
		fmt.Printf("Error dropping table: %v\n", dropError)
	}

	dropError = db.Migrator().DropTable(&AdvancedNumberFactors{})
	if dropError != nil {
		fmt.Printf("Error dropping table: %v\n", dropError)
	}

	// Migrate the schemas
	dbCreateError := db.AutoMigrate(&AdvancedNumberInformation{})
	if dbCreateError != nil {
		fmt.Printf("Error creating table: %v\n", dbCreateError)
	}

	dbCreateError = db.AutoMigrate(&AdvancedNumberFactors{})
	if dbCreateError != nil {
		fmt.Printf("Error creating table: %v\n", dbCreateError)
	}

	// close connection
	closeErr := liberdatabase.CloseConnection(db)
	if closeErr != nil {
		fmt.Printf("Error closing database connection: %v\n", closeErr)
	}

	return
}
