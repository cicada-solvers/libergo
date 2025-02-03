package main

import (
	"config"
	"context"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"liberdatabase"
	"math/big"
	"os"
	"sequences"
	"strings"
)

func main() {
	// Display a big warning outlined in red
	fmt.Println("\033[31;1m" + strings.Repeat("=", 80))
	fmt.Println("WARNING: This program is for use on puzzles only.")
	fmt.Println("Using this for other purposes may be illegal in your country.")
	fmt.Println(strings.Repeat("=", 80) + "\033[0m")

	// Define the combo flag
	comboFlag := flag.Bool("combo", false, "Find prime combos")
	flag.Parse()

	// Check if the number is provided as an argument
	if len(flag.Args()) < 1 {
		fmt.Println("Please provide a number to be factorized as an argument.")
		os.Exit(1)
	}

	// Read input number
	numberStr := flag.Arg(0)

	// Convert input to bigint
	number := new(big.Int)
	_, ok := number.SetString(numberStr, 10)
	if !ok {
		fmt.Println("Invalid number format.")
		os.Exit(1)
	}

	if number.Cmp(big.NewInt(1)) == -1 || number.Cmp(big.NewInt(1)) == 0 {
		fmt.Printf("%s : 1,%s\n", numberStr, numberStr)
		return
	}

	if sequences.IsPrime(number) {
		// You don't need to factorize a prime number
		fmt.Printf("%s : 1,%s\n", numberStr, numberStr)
		return
	}

	// Load database configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Establish database connection
	connStr := cfg.GeneralConnectionString
	db, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer func(db *pgx.Conn, ctx context.Context) {
		err := db.Close(ctx)
		if err != nil {
			fmt.Printf("Error closing connection: %v\n", err)
		}
	}(db, context.Background())

	// The mainId is the number being factorized
	mainId := uuid.New().String()

	if *comboFlag {
		// Find prime combos
		findCombos(db, mainId, number)

		// Output prime combos
		output := strings.Builder{}
		firstTime := true

		// Initialize the last sequence number
		var lastSeqNumber int64 = 0

		// Loop to get prime combos until nil is returned
		for {
			combo, err := liberdatabase.GetPrimeCombosByMainID(db, mainId, lastSeqNumber)
			if err != nil {
				fmt.Printf("Error getting prime combos: %v\n", err)
				os.Exit(1)
			}
			if combo == nil {
				break
			}

			// Update the last sequence number
			lastSeqNumber = combo.SeqNumber

			if !firstTime {
				output.WriteString(",")
			}

			// Append prime combo to output
			output.WriteString(fmt.Sprintf("%s : (%s,%s)", numberStr, combo.ValueP, combo.ValueQ))

			firstTime = false
		}

		fmt.Println(numberStr, ":", output.String())

		removeErr := liberdatabase.RemovePrimeCombosByMainID(db, mainId)
		if removeErr != nil {
			return
		}
	} else {
		// Perform factorization
		factorize(db, mainId, number, 0)

		// Output prime factors
		output := strings.Builder{}
		firstTime := true

		// Initialize the last sequence number
		var lastSeqNumber int64 = 0

		// Loop to get factors until nil is returned
		for {
			factor, err := liberdatabase.GetFactorsByMainID(db, mainId, lastSeqNumber)
			if err != nil {
				fmt.Printf("Error getting factors: %v\n", err)
				os.Exit(1)
			}
			if factor == nil {
				break
			}

			// Update the last sequence number
			lastSeqNumber = factor.SeqNumber

			if !firstTime {
				output.WriteString(",")
			}

			// Append factor to output
			output.WriteString(factor.Factor)

			firstTime = false
		}

		fmt.Println(numberStr, ":", output.String())

		removeErr := liberdatabase.RemoveFactorsByMainID(db, mainId)
		if removeErr != nil {
			return
		}
	}
}
