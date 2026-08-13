package main

import (
	trade "currency_converter/types"
	util "currency_converter/utils"
	"fmt"
	"os"
)

const ErrorMessage = "An error occurred: %v\n"

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("No input provided, use -help flag to see command to run")
		return
	}

	if os.Args[1] == "-help" {
		fmt.Println(`
		You can convert from one currency to another by providing the base currency first
		and the quote currency last with price of the base in quote. 
		
		Example: 
		go run . BTC USDT=20 price=1000
		`)
		return
	}

	if len(os.Args) < 4 {
		fmt.Println("Usage: BTC=1 USDT PRICE=20000")
		return
	}

	base, baseAmount, err := util.ExtractInputs(os.Args[1])

	if err != nil {
		fmt.Printf(ErrorMessage, err)
		return
	}

	quote, quoteAmount, err := util.ExtractInputs(os.Args[2])

	if err != nil {
		fmt.Printf(ErrorMessage, err)
		return
	}

	_, marketPrice, err := util.ExtractInputs(os.Args[3])

	if err != nil {
		fmt.Printf(ErrorMessage, err)
		return
	}

	tradeQuote := trade.NewTradeQuote(base, quote, baseAmount, quoteAmount, marketPrice)

	result, err := tradeQuote.Execute()

	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println("Your trade:", result)
}
