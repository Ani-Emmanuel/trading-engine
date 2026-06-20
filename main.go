package main

import (
	"fmt"
	"os"
)

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

	base, baseAmount, err := extractInputs(os.Args[1])

	if err != nil {
		fmt.Printf(ErrorMessage, err)
		return
	}

	quote, quoteAmount, err := extractInputs(os.Args[2])

	if err != nil {
		fmt.Printf(ErrorMessage, err)
		return
	}

	_, marketPrice, err := extractInputs(os.Args[3])

	if err != nil {
		fmt.Printf(ErrorMessage, err)
		return
	}

	tradeQuote := NewTradeQuote(base, quote, baseAmount, quoteAmount, marketPrice)

	result, err := tradeQuote.Execute()

	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println("Your trade:", result)
}
