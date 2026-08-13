package types

import "fmt"

type TradeResult struct {
	CurrencyCode string
	Quantity     float64
}

func (r TradeResult) String() string {
	return fmt.Sprintf("%.2f %s", r.Quantity, r.CurrencyCode)
}
