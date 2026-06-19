package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

	result, err := trade(base, baseAmount, quote, quoteAmount, marketPrice)

	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println("Your trade", result)
}

func trade(base string, baseAmount float64, quote string, quoteAmount float64, price float64) (map[string]float64, error) {
	result := make(map[string]float64)

	if baseAmount > 0 && quoteAmount > 0 {
		return nil, fmt.Errorf("Invalid input. You need only to specify either the base or quote price basePrice: %v, quotePrice: %v\n", baseAmount, quoteAmount)
	}

	if baseAmount > 0 {
		amount := baseAmount * price
		result[quote] = amount
	}

	if quoteAmount > 0 {
		amount := quoteAmount / price
		result[base] = amount
	}

	return result, nil
}

func extractAmount(amount string) (float64, error) {
	if amount == "" {
		return 0, fmt.Errorf("input cannot be empty")
	}

	parsedAmount, err := strconv.ParseFloat(amount, 64)

	if err != nil {
		return 0, fmt.Errorf("invalid input: %s", amount)
	}

	if parsedAmount <= 0 {
		return 0, fmt.Errorf("input %s must be greater than zero", amount)
	}

	return parsedAmount, nil
}

func extractInputs(input string) (string, float64, error) {
	result := strings.SplitN(input, "=", 2)

	if len(result) > 1 {
		code := strings.ToUpper(strings.TrimSpace(result[0]))
		amount, err := extractAmount(strings.TrimSpace(result[1]))
		return code, amount, err
	}

	return strings.ToUpper(strings.TrimSpace(result[0])), 0, nil

}
