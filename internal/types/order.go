package types

import "time"

type Order struct {
	ID        string
	OwnerId   string
	Qty       int
	Price     int
	Code      string
	Side      Side
	CreatedAt time.Time
}
