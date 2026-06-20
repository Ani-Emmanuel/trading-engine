package main

import (
	"fmt"
	"strconv"
	"strings"
)

const ErrorMessage = "An error occurred: %v\n"

type TradeQuote struct {
	Base          string
	Quote         string
	BaseQuantity  float64
	QuoteQuantity float64
	Price         float64
}

type TradeResult struct {
	CurrencyCode string
	Quantity     float64
}

type ValidationError struct {
	Message string
}

func (t TradeQuote) Execute() (TradeResult, error) {
	if err := t.Validate(); err != nil {
		return TradeResult{}, err
	}

	if t.BaseQuantity > 0 {
		return TradeResult{
			CurrencyCode: t.Quote,
			Quantity:     t.BaseQuantity * t.Price,
		}, nil
	}

	return TradeResult{
		CurrencyCode: t.Base,
		Quantity:     t.QuoteQuantity / t.Price,
	}, nil
}

func (t TradeQuote) Validate() error {
	if t.BaseQuantity > 0 && t.QuoteQuantity > 0 {
		return ValidationError{Message: fmt.Sprintf("Both quantities provided: base=%v quote=%v",
			t.BaseQuantity,
			t.QuoteQuantity,
		)}
	}

	if t.BaseQuantity == 0 && t.QuoteQuantity == 0 {
		return ValidationError{
			Message: fmt.Sprintf("No quantity specified: base=%v quote=%v", t.BaseQuantity,
				t.QuoteQuantity),
		}
	}

	if t.Price <= 0 {
		return ValidationError{
			Message: fmt.Sprintf("Price must be greater than zero: price=%v", t.Price),
		}
	}

	return nil
}

func (r TradeResult) String() string {
	return fmt.Sprintf("%.2f %s", r.Quantity, r.CurrencyCode)
}

func (e ValidationError) Error() string {
	return e.Message
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

func NewTradeQuote(base string, quote string, baseQuantity float64, quoteQuantity float64, price float64) TradeQuote {
	return TradeQuote{
		Base:          base,
		Quote:         quote,
		BaseQuantity:  baseQuantity,
		QuoteQuantity: quoteQuantity,
		Price:         price,
	}
}
