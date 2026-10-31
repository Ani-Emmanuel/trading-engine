package types

import "time"

type Order struct {
	ID        string
	OwnerID   string
	Qty       int
	Price     int
	Code      string
	Side      Side
	CreatedAt time.Time
}
