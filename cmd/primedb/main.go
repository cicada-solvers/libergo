package main

import (
	"liberdatabase"
	"math"
	"math/big"
	"sequences"
)

func main() {
	_, _ = liberdatabase.InitTables()
	conn, _ := liberdatabase.InitConnection()
	one := big.NewInt(1)
	number := big.NewInt(2)
	maxNumber := big.NewInt(math.MaxInt32)
	for number.Cmp(maxNumber) <= 0 {
		if sequences.IsPrime(number) {
			liberdatabase.AddPrimeValue(conn, number.Int64())
		}
		number = number.Add(number, one)
	}

	_ = liberdatabase.CloseConnection(conn)
}
