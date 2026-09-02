package orderbook

import (
	"time"

	"github.com/ani-emmanuel/trading-engine/types"

	"github.com/google/uuid"
)

type OrderBook struct {
	BuyBook  map[int][]types.Order
	SellBook map[int][]types.Order
}

func (ob *OrderBook) AddToOrderBook(userId, code string, side types.Side, qty, price int) error {
	book := types.Order{
		ID:        uuid.NewString(),
		OwnerId:   userId,
		Qty:       qty,
		Code:      code,
		Price:     price,
		Side:      side,
		CreatedAt: time.Now(),
	}

	orderBookMap, err := determineOrderBookBySide(ob, side)

	if err != nil {
		return err
	}

	orderBookMap[price] = append(orderBookMap[price], book)

	return nil
}

func (ob *OrderBook) RemoveFromOrderBook(side types.Side, price, index int) error {
	orderBook, err := determineOrderBookBySide(ob, side)

	if err != nil {
		return err
	}

	modifyOrderBook(orderBook, price, index)

	return nil
}

func (ob *OrderBook) UpdateOrderBookRecord(side types.Side, qty, price, index int) error {
	orderBook, err := determineOrderBookBySide(ob, side)

	if err != nil {
		return err
	}

	return updateOrderBookQuantity(orderBook, qty, price, index)
}

func (ob *OrderBook) GetOrderBookRecord(side types.Side, price, index int) (types.Order, error) {
	orders, err := determineOrderBookBySide(ob, side)

	if err != nil {
		return types.Order{}, err
	}

	return fetchOrderBookRecord(orders, price, index)
}

func (ob *OrderBook) GetBook(side types.Side) (map[int][]types.Order, error) {
	return determineOrderBookBySide(ob, side)
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		BuyBook:  make(map[int][]types.Order),
		SellBook: make(map[int][]types.Order),
	}
}

func modifyOrderBook(orderBook map[int][]types.Order, price, index int) map[int][]types.Order {
	orders := orderBook[price]

	if index < 0 || index >= len(orders) {
		return orderBook
	}

	orderBook[price] = append(
		orders[:index],
		orders[index+1:]...,
	)

	if len(orderBook[price]) == 0 {
		delete(orderBook, price)
	}

	return orderBook
}

func updateOrderBookQuantity(orderBook map[int][]types.Order, qty, price, index int) error {
	orders := orderBook[price]

	if index < 0 || index >= len(orders) {
		return types.ValidationError{Message: "Order record does not exist"}
	}

	orders[index].Qty = qty

	return nil
}

func fetchOrderBookRecord(orderBook map[int][]types.Order, price, index int) (types.Order, error) {
	orders := orderBook[price]

	if index < 0 || index >= len(orders) {
		return types.Order{}, types.ValidationError{Message: "Order record does not exist"}
	}

	return orders[index], nil
}

func determineOrderBookBySide(orderBook *OrderBook, side types.Side) (map[int][]types.Order, error) {
	switch side {
	case types.Buy:
		return orderBook.BuyBook, nil
	case types.Sell:
		return orderBook.SellBook, nil
	default:
		return nil, types.ValidationError{Message: "Invalid side"}
	}
}
