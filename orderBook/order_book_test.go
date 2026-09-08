package orderbook

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ani-emmanuel/trading-engine/types"
)

func TestAddToOrderBook(t *testing.T) {

	t.Run("adds a sell order", func(t *testing.T) {
		ob := NewOrderBook()
		requireNoError(t, ob.AddToOrderBook("user-ID-1", "BTC", types.Sell, 10, 100000))

		orders := ob.SellBook[100000]
		assertPriceLevelCount(t, ob.SellBook, 1)
		assertOrderCount(t, orders, 1)
		assertOrder(t, orders[0], "user-ID-1", "BTC", types.Sell, 10, 100000)
	})

	t.Run("adds multiple orders at the same price", func(t *testing.T) {
		ob := NewOrderBook()

		for x := range 3 {
			userID := fmt.Sprintf("user-ID-%d", x)
			err := ob.AddToOrderBook(userID, "BTC", types.Sell, 5, 100000)
			requireNoError(t, err)
		}

		orders := ob.SellBook[100000]
		assertPriceLevelCount(t, ob.SellBook, 1)
		assertOrderCount(t, orders, 3)
		assertOrder(t, orders[0], "user-ID-0", "BTC", types.Sell, 5, 100000)
		assertOrder(t, orders[1], "user-ID-1", "BTC", types.Sell, 5, 100000)
		assertOrder(t, orders[2], "user-ID-2", "BTC", types.Sell, 5, 100000)
	})

	t.Run("adds multiple orders at different prices", func(t *testing.T) {
		ob := NewOrderBook()

		for v := range 3 {
			userID := fmt.Sprintf("user-ID-%d", v)
			err := ob.AddToOrderBook(userID, "BTC", types.Sell, 10, 3000)
			requireNoError(t, err)
		}

		for v := range 2 {
			userID := fmt.Sprintf("user-ID-%d", v)
			err := ob.AddToOrderBook(userID, "BTC", types.Buy, 10, 2000)
			requireNoError(t, err)
		}

		sellOrders := ob.SellBook[3000]
		assertPriceLevelCount(t, ob.SellBook, 1)
		assertOrderCount(t, sellOrders, 3)
		assertOrder(t, sellOrders[0], "user-ID-0", "BTC", types.Sell, 10, 3000)

		buyOrders := ob.BuyBook[2000]
		assertPriceLevelCount(t, ob.BuyBook, 1)
		assertOrderCount(t, buyOrders, 2)
		assertOrder(t, buyOrders[0], "user-ID-0", "BTC", types.Buy, 10, 2000)
	})
}

func TestRemoveFromOrderBook(t *testing.T) {
	t.Run("removes from sell order", func(t *testing.T) {
		ob := NewOrderBook()

		for x := range 3 {
			value := x + 1

			userID := fmt.Sprintf("user-ID-%d", value)

			err := ob.AddToOrderBook(userID, "BTC", types.Sell, 5, 100)
			requireNoError(t, err)
		}

		err := ob.RemoveFromOrderBook(types.Sell, 100, 0)
		requireNoError(t, err)

		orders := ob.SellBook[100]

		assertPriceLevelCount(t, ob.SellBook, 1)
		assertOrderCount(t, orders, 2)
		assertOrder(t, orders[0], "user-ID-2", "BTC", types.Sell, 5, 100)
		assertOrder(t, orders[1], "user-ID-3", "BTC", types.Sell, 5, 100)
	})

	t.Run("removes order book", func(t *testing.T) {
		ob := NewOrderBook()

		requireNoError(t, ob.AddToOrderBook("user-ID-1", "BTC", types.Sell, 5, 100))

		err := ob.RemoveFromOrderBook(types.Sell, 100, 0)
		requireNoError(t, err)

		_, exists := ob.SellBook[100]

		if exists {
			t.Fatalf("expected price level to be removed, but it still exists")
		}
	})

	t.Run("does not remove when index does not exist", func(t *testing.T) {
		ob := NewOrderBook()
		requireNoError(t, ob.AddToOrderBook("user-ID-1", "BTC", types.Sell, 5, 100))

		err := ob.RemoveFromOrderBook(types.Sell, 100, 1)
		requireNoError(t, err)

		_, exists := ob.SellBook[100]

		if !exists {
			t.Fatalf("expected price level to not be removed, but it was removed")
		}

		assertPriceLevelCount(t, ob.SellBook, 1)
		assertOrderCount(t, ob.SellBook[100], 1)
		assertOrder(t, ob.SellBook[100][0], "user-ID-1", "BTC", types.Sell, 5, 100)
	})
}

func TestUpdateOrderBookRecord(t *testing.T) {
	t.Run("updates order quantity", func(t *testing.T) {
		ob := NewOrderBook()

		requireNoError(t, ob.AddToOrderBook("user-ID-1", "BTC", types.Sell, 5, 100))
		requireNoError(t, ob.AddToOrderBook("user-ID-2", "BTC", types.Sell, 5, 100))

		err := ob.UpdateOrderBookRecord(types.Sell, 3, 100, 1)

		requireNoError(t, err)

		order := ob.SellBook[100][1]

		assertOrder(t, order, "user-ID-2", "BTC", types.Sell, 3, 100)
	})

	t.Run("returns error if index does not exist", func(t *testing.T) {
		ob := NewOrderBook()

		requireNoError(t, ob.AddToOrderBook("user-ID-1", "BTC", types.Sell, 5, 100))

		err := ob.UpdateOrderBookRecord(types.Sell, 3, 100, 1)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var validationError types.ValidationError

		if !errors.As(err, &validationError) {
			t.Fatalf("expected ValidationError, got %T", err)
		}

		if validationError.Message != "Order record does not exist" {
			t.Errorf("expected message Order record does not exist, got %q", validationError.Message)
		}
	})
}

func TestGetOrderBookRecord(t *testing.T) {
	t.Run("returns order at given index", func(t *testing.T) {
		ob := NewOrderBook()

		requireNoError(t, ob.AddToOrderBook("user-ID-1", "BTC", types.Sell, 5, 100))
		requireNoError(t, ob.AddToOrderBook("user-ID-2", "BTC", types.Sell, 5, 100))

		order, err := ob.GetOrderBookRecord(types.Sell, 100, 1)

		if err != nil {
			t.Fatalf("expected no err, got %q", err)
		}

		assertOrder(t, order, "user-ID-2", "BTC", types.Sell, 5, 100)
	})

	t.Run("returns error when index does not exist", func(t *testing.T) {
		ob := NewOrderBook()

		requireNoError(t, ob.AddToOrderBook("user-ID-1", "BTC", types.Sell, 5, 100))

		_, err := ob.GetOrderBookRecord(types.Sell, 100, 1)

		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		var validationError types.ValidationError

		if !errors.As(err, &validationError) {
			t.Fatalf("expected ValidationError, got %T", err)
		}

		if validationError.Message != "Order record does not exist" {
			t.Errorf("expected message Order record does not exist, got %q", validationError.Message)
		}
	})
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func assertPriceLevelCount(t *testing.T, book map[int][]types.Order, expected int) {
	t.Helper()
	if len(book) != expected {
		t.Fatalf("expected %d price level, got %d", expected, len(book))
	}
}

func assertOrderCount(t *testing.T, orders []types.Order, expected int) {
	t.Helper()
	if len(orders) != expected {
		t.Fatalf("expected %d order, got %d", expected, len(orders))
	}
}

func assertOrder(t *testing.T, order types.Order, ownerID, code string, side types.Side, qty, price int) {
	t.Helper()

	if order.OwnerId != ownerID {
		t.Fatalf("expected owner id %s, got %s", ownerID, order.OwnerId)
	}

	if order.Qty != qty {
		t.Fatalf("expected quantity %d, got %d", qty, order.Qty)
	}

	if order.Price != price {
		t.Fatalf("expected price %d, got %d", price, order.Price)
	}

	if order.Code != code {
		t.Fatalf("expected code %s, got %s", code, order.Code)
	}

	if order.Side != side {
		t.Fatalf("expected side %v, got %v", side, order.Side)
	}
}
