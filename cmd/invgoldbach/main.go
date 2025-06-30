package main

import (
	"fmt"
	"gorm.io/gorm"
	"liberdatabase"
	"numeric"
)

func main() {
	// Initialize database
	_, _ = liberdatabase.InitTables()
	conn, _ := liberdatabase.InitConnection()
	defer func(db *gorm.DB) {
		err := liberdatabase.CloseConnection(db)
		if err != nil {
			fmt.Println("Error closing database connection:", err)
		}
	}(conn)

	largestNumber := int64(2147483647) + int64(2147483647)

	for i := int64(4); i <= largestNumber; i++ {
		if numeric.IsNumberEven(i) {
			gbp := numeric.NewGoldbachPairs()
			primeList := liberdatabase.GetPrimeListLessThanValue(conn, i)

			solveError := gbp.SolveForNumber(i, &primeList)
			if solveError != nil {
				fmt.Println("Error solving for number:", i, solveError)
				fmt.Println("----------------------------------------")
			} else {
				pairCount := len(gbp.GetGoldbachPairs())
				fmt.Println("Number:", i, pairCount)
				fmt.Println("----------------------------------------")

				// Now we need to add them to the database
				goldbachNumber := liberdatabase.AddGoldbachNumber(conn, i, true)
				for _, pair := range gbp.GetGoldbachPairs() {
					liberdatabase.AddGoldbachAddend(conn, goldbachNumber.Id, pair.AddendOne, pair.AddendTwo)
				}
			}

		}
	}
}
