package main

import (
	"fmt"
	"liberdatabase"
	"math"
	"math/big"
	"sequences"
)

func main() {
	_, _ = liberdatabase.InitTables()
	conn, _ := liberdatabase.InitConnection()
	for i := int64(2); i < math.MaxInt64; i++ {
		bigI := big.NewInt(int64(i))
		if sequences.IsPrime(bigI) {
			liberdatabase.AddPrimeValue(conn, i)
			fmt.Println(i)
		}
	}

	_ = liberdatabase.CloseConnection(conn)
}
