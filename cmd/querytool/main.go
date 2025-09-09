package main

import (
	"flag"
	"fmt"
	"liberdatabase"
	"strconv"
)

func main() {
	view := flag.String("view", "runelength", "The view to call")
	params := flag.String("params", "3", "The parameters to pass to the view")

	flag.Parse()

	if *view == "" {
		flag.Usage()
		return
	}

	dbconn, _ := liberdatabase.InitConnection()

	if *view == "runelength" {
		length := int64(0)
		if *params != "" {
			length, _ = strconv.ParseInt(*params, 10, 64)
		}

		list := liberdatabase.GetDictionaryWordsByRuneLength(dbconn, int(length))

		for _, word := range list {
			fmt.Printf("%s\n", word)
		}
	}

	_ = liberdatabase.CloseConnection(dbconn)
	return
}
