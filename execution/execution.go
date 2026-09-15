package execution

import (
	"time"

	"github.com/google/uuid"
)

type Execution struct {
	ID           string
	MakerOrderID string
	TakerOrderID string
	MakerUserID  string
	TakerUserID  string
	Code         string
	Price        int
	Quantity     int
	CreatedAt    time.Time
}

func NewExecution(
	makerOrderID, takerOrderID,
	makerUserID, takerUserID,
	code string,
	price, quantity int,
) Execution {
	return Execution{
		ID:           uuid.NewString(),
		MakerOrderID: makerOrderID,
		TakerOrderID: takerOrderID,
		MakerUserID:  makerUserID,
		TakerUserID:  takerUserID,
		Code:         code,
		Price:        price,
		Quantity:     quantity,
		CreatedAt:    time.Now(),
	}
}
