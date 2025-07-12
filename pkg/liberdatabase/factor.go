package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"
)

// Factor represents a database model mapped to the "factors" table.
// It contains fields for ID, Factor, MainId, and SeqNumber for storage and retrieval operations.
type Factor struct {
	ID        string `gorm:"column:id"`
	Factor    string `gorm:"column:factor"`
	MainId    string `gorm:"column:mainid"`
	SeqNumber int64  `gorm:"column:seqnumber"`
}

// TableName specifies the database table name "factors" for the Factor model.
func (Factor) TableName() string {
	return "factors"
}

// GetFactorsByMainID retrieves all factors from the factors table based on the mainid.
func GetFactorsByMainID(db *gorm.DB, mainId string, seqNumber int64) *Factor {
	var count int64
	counterResult := db.Model(&Factor{}).Where("mainid = ? AND seqnumber > ?", mainId, seqNumber).Count(&count)
	if counterResult.Error != nil {
		fmt.Printf("error counting factors: %v", counterResult.Error)
	}

	if count == 0 {
		return nil
	}

	var factor Factor
	result := db.Where("mainid = ? AND seqnumber > ?", mainId, seqNumber).Order("seqnumber ASC").Limit(1).First(&factor)
	if result.RowsAffected == 0 {
		return nil
	}
	if result.Error != nil {
		fmt.Printf("error querying factors: %v\n", result.Error)
	}
	return &factor
}

// GetMaxSeqNumberByMainID retrieves the maximum SeqNumber from the factors table based on the mainid.
func GetMaxSeqNumberByMainID(db *gorm.DB, mainId string) Factor {
	var factor Factor
	result := db.Where("mainid = ?", mainId).Order("seqnumber DESC").Limit(1).First(&factor)
	if result.Error != nil {
		fmt.Printf("error querying factors: %v\n", result.Error)
	}
	return factor
}

// InsertFactor inserts a factor into the database
func InsertFactor(db *gorm.DB, factor Factor) {
	result := db.Create(&factor)
	if result.Error != nil {
		fmt.Printf("error inserting factor: %v\n", result.Error)
	}
}

// RemoveFactorByID removes a factor from the factors table based on the given id.
func RemoveFactorByID(db *gorm.DB, id string) {
	result := db.Delete(&Factor{}, "id = ?", id)
	if result.Error != nil {
		fmt.Println("error deleting factors: %v\n", result.Error)
	}
}

// RemoveFactorsByMainID removes all factors from the factors table based on the given mainId.
func RemoveFactorsByMainID(db *gorm.DB, mainId string) {
	result := db.Delete(&Factor{}, "mainid = ?", mainId)
	if result.Error != nil {
		fmt.Println("error deleting factors: %v\n", result.Error)
	}
}
