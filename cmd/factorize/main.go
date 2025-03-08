package main

import (
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"liberdatabase"
	"math/big"
	"os"
	"strings"
	"time"
	"titler"
)

func main() {
	titler.PrintTitle("Factorize")
	startTime := time.Now()

	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Please provide a number to be factorized as an argument.")
		os.Exit(1)
	}

	numberStr := flag.Arg(0)
	number := new(big.Int)
	number, ok := number.SetString(numberStr, 10)
	if !ok {
		fmt.Println("Invalid number format.")
		os.Exit(1)
	}

	if number.Cmp(big.NewInt(1)) == -1 || number.Cmp(big.NewInt(1)) == 0 {
		fmt.Printf("%s : 1,%s\n", numberStr, numberStr)
		return
	}

	if number.ProbablyPrime(20) {
		fmt.Printf("%s : 1,%s\n", numberStr, numberStr)
		return
	}

	createErr := liberdatabase.InitSQLiteTables()
	if createErr != nil {
		fmt.Printf("Error initializing SQLite tables: %v\n", createErr)
		os.Exit(1)
	}

	db, connError := liberdatabase.InitConnection()
	if connError != nil {
		fmt.Printf("Error initializing database connection: %v\n", connError)
		os.Exit(1)
	}

	mainId := uuid.New().String()

	fmt.Printf("Processing %s (%d bits)...\n", numberStr, number.BitLen())

	p := tea.NewProgram(model{counter: big.NewInt(0)})
	go func() {
		factorize(db, mainId, number, 0, p)
		p.Send(tea.Quit())
	}()

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting Bubble Tea program: %v\n", err)
		os.Exit(1)
	}

	output := strings.Builder{}
	firstTime := true
	var lastSeqNumber int64 = 0

	for {
		factor := liberdatabase.GetFactorsByMainID(db, mainId, lastSeqNumber)
		if factor == nil {
			break
		}

		lastSeqNumber = factor.SeqNumber

		if !firstTime {
			output.WriteString(",")
		}

		output.WriteString(factor.Factor)
		firstTime = false
	}

	fmt.Println(numberStr, ":", output.String())
	liberdatabase.RemoveFactorsByMainID(db, mainId)

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("Execution time: %v\n", duration)
}
