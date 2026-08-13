package types

import (
	"fmt"
)

type TradeQuote struct {
	Base          string
	Quote         string
	BaseQuantity  float64
	QuoteQuantity float64
	Price         float64
}

type TradeQuoteInterface interface {
	Execute() (TradeResult, error)
	Validate() error
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

func NewTradeQuote(base string, quote string, baseQuantity float64, quoteQuantity float64, price float64) TradeQuote {
	return TradeQuote{
		Base:          base,
		Quote:         quote,
		BaseQuantity:  baseQuantity,
		QuoteQuantity: quoteQuantity,
		Price:         price,
	}
}
