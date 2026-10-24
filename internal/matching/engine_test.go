package matching

import (
	"fmt"
	"testing"

	"trading-engine/internal/execution"
	orderbook "trading-engine/internal/orderBook"
	"trading-engine/internal/types"
)

func TestMatchEngine(t *testing.T) {
	t.Run("creates resting order when book is empty", testCreatesRestingOrder)
	t.Run("fully matches maker and taker", testFullyMatchesMakerAndTaker)
	t.Run("matches multiple makers at the same price in FIFO order", testMatchesMultipleMakers)
	t.Run("partially fills maker", testPartiallyFillsMaker)
	t.Run("partially fills taker", testPartiallyFillsTaker)
	t.Run("matches best price first", testPricePriority)
	t.Run("reset if the order did not match", testWhenOrderDoNotMatch)
	t.Run("partially matches order and rests remaining taker quantity", testPartialMatchAndRestRemainingQuantity)
	t.Run("create a resting order if it can't match", testOrderDoesNotCrossPrice)
}

func testCreatesRestingOrder(t *testing.T) {
	ob := orderbook.NewOrderBook()

	engine := Engine{
		OrderBook: ob,
	}

	order := types.Order{
		OwnerId: "user-1",
		Qty:     2,
		Price:   1000,
		Code:    "BTC",
		Side:    types.Sell,
	}

	_, err := engine.Match(order)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	book, err := engine.OrderBook.GetBook(order.Side)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	orders := book[order.Price]

	if len(orders) != 1 {
		t.Fatalf("expected 1 order, got %d", len(orders))
	}

	if orders[0].OwnerId != order.OwnerId {
		t.Errorf("expected owner ID %q, got %q", order.OwnerId, orders[0].OwnerId)
	}

	if orders[0].Qty != order.Qty {
		t.Errorf("expected quantity %d, got %d", order.Qty, orders[0].Qty)
	}

	if orders[0].Price != order.Price {
		t.Errorf("expected price %d, got %d", order.Price, orders[0].Price)
	}

	if orders[0].Side != order.Side {
		t.Errorf("expected side %q, got %q", order.Side, orders[0].Side)
	}

}

func testFullyMatchesMakerAndTaker(t *testing.T) {
	ob := orderbook.NewOrderBook()

	engine := Engine{
		OrderBook: ob,
	}

	makerOrder := types.Order{
		OwnerId: "user-1",
		Qty:     2,
		Price:   1000,
		Code:    "BTC",
		Side:    types.Buy,
	}

	takerOrder := types.Order{
		OwnerId: "user-2",
		Qty:     2,
		Price:   1000,
		Code:    "BTC",
		Side:    types.Sell,
	}

	_, err := engine.Match(makerOrder)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	executions, err := engine.Match(takerOrder)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(executions) != 1 {
		t.Fatalf("expected 1 expression, got %d", len(executions))
	}

	if executions[0].Code != "BTC" {
		t.Fatalf("expected BTC, got %s", executions[0].Code)
	}

	book, err := engine.OrderBook.GetBook(types.Buy)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(book[1000]) != 0 {
		t.Fatalf("expected maker order to be removed")
	}
}

func testMatchesMultipleMakers(t *testing.T) {
	ob := orderbook.NewOrderBook()

	engine := Engine{
		OrderBook: ob,
	}

	sellOrder := types.Order{
		OwnerId: "user-1",
		Qty:     3,
		Price:   1000,
		Code:    "BTC",
		Side:    types.Sell,
	}

	for v := range 2 {
		if err := engine.OrderBook.AddToOrderBook(fmt.Sprintf("user-%d", v), "BTC", types.Buy, 1, 1000); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	}

	executions, err := engine.Match(sellOrder)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(executions) != 2 {
		t.Fatalf("expected 2, got %d", len(executions))
	}

	record, _ := engine.OrderBook.GetBook(types.Sell)

	data := record[1000]

	if len(data) != 1 {
		t.Fatalf("expected 1 record, got %d", len(data))
	}

	if data[0].OwnerId != "user-1" {
		t.Fatalf("expected user-1, got %s", data[0].OwnerId)
	}

	if data[0].Side != types.Sell {
		t.Fatalf("expected Sell, got %s", types.Sell)
	}

	if data[0].Qty != 1 {
		t.Fatalf("expected remaining quantity 1, got %d", data[0].Qty)
	}

	for i, data := range executions {
		if data.MakerUserID != fmt.Sprintf("user-%d", i) {
			t.Fatalf("expected %s, got %s", fmt.Sprintf("user-%d", i), data.MakerUserID)
		}
	}
}

func testPartiallyFillsMaker(t *testing.T) {
	ob := orderbook.NewOrderBook()

	engine := Engine{
		OrderBook: ob,
	}

	taker := types.Order{
		OwnerId: "user-1",
		Qty:     2,
		Price:   1000,
		Code:    "BTC",
		Side:    types.Buy,
	}

	if err := engine.OrderBook.AddToOrderBook("user-1", "BTC", types.Sell, 3, 1000); err != nil {
		t.Fatalf("expected error to be nil, got %v", err)
	}

	executions, err := engine.Match(taker)

	if err != nil {
		t.Fatalf("expected error to nil, got %v", err)
	}

	if len(executions) != 1 {
		t.Fatalf("expected 1 execution, got %d", len(executions))
	}

	orders, err := engine.OrderBook.GetBook(types.Sell)

	if err != nil {
		t.Fatalf("expected error to nil, got %v", err)
	}

	order := orders[1000]

	if len(order) != 1 {
		t.Fatalf("expected order to be 1, got %d", len(order))
	}

	if order[0].Price != 1000 {
		t.Fatalf("expected order price be 1000, got %d", order[0].Price)
	}

	if order[0].Qty != 1 {
		t.Fatalf("expected order qty be 1, got %d", order[0].Qty)
	}
}

func testPartiallyFillsTaker(t *testing.T) {
	ob := orderbook.NewOrderBook()

	engine := Engine{
		OrderBook: ob,
	}

	taker := types.Order{
		OwnerId: "user-1",
		Qty:     5,
		Price:   1000,
		Code:    "BTC",
		Side:    types.Buy,
	}

	if err := engine.OrderBook.AddToOrderBook("user-1", "BTC", types.Sell, 3, 1000); err != nil {
		t.Fatalf("expected error to be nil, got %v", err)
	}

	executions, err := engine.Match(taker)

	if err != nil {
		t.Fatalf("expected error to nil, got %v", err)
	}

	if len(executions) != 1 {
		t.Fatalf("expected 1 execution, got %d", len(executions))
	}

	buyOrders, err := engine.OrderBook.GetBook(types.Buy)

	if err != nil {
		t.Fatalf("expected error to nil, got %v", err)
	}

	order := buyOrders[1000]

	if len(order) != 1 {
		t.Fatalf("expected order to be 1, got %d", len(order))
	}

	if order[0].Price != 1000 {
		t.Fatalf("expected order price be 1000, got %d", order[0].Price)
	}

	if order[0].Qty != 2 {
		t.Fatalf("expected order qty be 1, got %d", order[0].Qty)
	}

	SellOrders, err := engine.OrderBook.GetBook(types.Sell)

	if err != nil {
		t.Fatalf("expected error to nil, got %v", err)
	}

	if _, exists := SellOrders[1000]; exists {
		t.Fatalf("expected sell price level 1000 to be removed")
	}
}

func testPricePriority(t *testing.T) {
	ob := orderbook.NewOrderBook()

	engine := Engine{
		OrderBook: ob,
	}

	price := 1100
	for v := range 3 {
		if err := engine.OrderBook.AddToOrderBook(fmt.Sprintf("user-%d", v), "BTC", types.Sell, 1, price); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		price -= 50
	}

	taker := types.Order{
		OwnerId: "user-1",
		Qty:     2,
		Price:   1100,
		Code:    "BTC",
		Side:    types.Buy,
	}

	executions, err := engine.Match(taker)

	if err != nil {
		t.Fatalf("expected error to nil, got %v", err)
	}

	if len(executions) != 2 {
		t.Fatalf("expected 2 executions, got %d", len(executions))
	}

	if executions[0].MakerUserID != "user-2" {
		t.Fatalf("expected owner id to be user-2, got %s", executions[0].MakerUserID)
	}

	if executions[0].Quantity != 1 {
		t.Fatalf("expected 1 quantity, got %d", executions[0].Quantity)
	}

	if executions[0].Price != 1000 {
		t.Fatalf("expected price to be 1000, got %d", executions[0].Price)
	}

	if executions[1].MakerUserID != "user-1" {
		t.Fatalf("expected owner id to be user-1, got %s", executions[1].MakerUserID)
	}

	if executions[1].Quantity != 1 {
		t.Fatalf("expected 1 quantity, got %d", executions[1].Quantity)
	}

	if executions[1].Price != 1050 {
		t.Fatalf("expected price to be 1050, got %d", executions[1].Price)
	}
}

func testWhenOrderDoNotMatch(t *testing.T) {
	ob := orderbook.NewOrderBook()
	engine := Engine{
		OrderBook: ob,
	}

	if err := engine.OrderBook.AddToOrderBook("user-1", "BTC", types.Sell, 1, 1100); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err := engine.Match(types.Order{
		OwnerId: "user-1",
		Qty:     5,
		Price:   1000,
		Code:    "BTC",
		Side:    types.Buy,
	})

	if err != nil {
		t.Fatalf("expected error to nil, got %v", err)
	}

	sellBook, err := engine.OrderBook.GetBook(types.Sell)

	if err != nil {
		t.Fatalf("expected error to nil, got %v", err)
	}

	buyBook, err := engine.OrderBook.GetBook(types.Buy)

	if err != nil {
		t.Fatalf("expected error to nil, got %v", err)
	}

	if len(sellBook[1100]) != 1 {
		t.Fatalf("expected 1 sell order, got %d", len(sellBook[1100]))
	}

	if len(buyBook[1000]) != 1 {
		t.Fatalf("expected 1 buy order, got %d", len(buyBook[1000]))
	}

	if sellBook[1100][0].Qty != 1 {
		t.Fatalf("expected 1 quantity, got %d", sellBook[1100][0].Qty)
	}

	if buyBook[1000][0].Qty != 5 {
		t.Fatalf("expected 5 quantity, got %d", buyBook[1000][0].Qty)
	}
}

func testPartialMatchAndRestRemainingQuantity(t *testing.T) {
	ob := orderbook.NewOrderBook()
	engine := Engine{
		OrderBook: ob,
	}

	price := 2000

	for v := range 2 {
		if err := engine.OrderBook.AddToOrderBook(fmt.Sprintf("user-%d", v), "BTC", types.Sell, (2 - v), price); err != nil {
			t.Fatalf("expected error to nil, got %v", err)
		}

		price -= 1000
	}

	execution, err := engine.Match(types.Order{
		OwnerId: "user-1",
		Qty:     3,
		Price:   1000,
		Code:    "BTC",
		Side:    types.Buy,
	})
	mustNoErr(t, err)
	assertExecution(t, execution, 1, 1, 1000)

	buyBook, err := engine.OrderBook.GetBook(types.Buy)
	mustNoErr(t, err)

	sellBook, err := engine.OrderBook.GetBook(types.Sell)
	mustNoErr(t, err)

	remainingBuyOrders, exists := buyBook[1000]
	if !exists {
		t.Fatal("expected buy price level 1000 to exist")
	}
	assertRestingOrders(t, remainingBuyOrders, 1, 2, 1000, "remainingBuyOrders")

	remainingSellOrders, exists := sellBook[2000]
	if !exists {
		t.Fatal("expected sell price level 2000 to exist")
	}
	assertRestingOrders(t, remainingSellOrders, 1, 2, 2000, "remainingSellOrders")
}

func testOrderDoesNotCrossPrice(t *testing.T) {
	ob := orderbook.NewOrderBook()
	engine := Engine{
		OrderBook: ob,
	}

	if err := engine.OrderBook.AddToOrderBook("user-1", "BTC", types.Sell, 1, 1100); err != nil {
		t.Fatalf("expected error to nil, got %v", err)
	}

	execution, err := engine.Match(types.Order{
		OwnerId: "user-2",
		Qty:     3,
		Price:   1000,
		Code:    "BTC",
		Side:    types.Buy,
	})

	mustNoErr(t, err)

	if len(execution) != 0 {
		t.Fatalf("expected 0 execution, got %d", len(execution))
	}

	buyBook, err := engine.OrderBook.GetBook(types.Buy)
	mustNoErr(t, err)

	sellBook, err := engine.OrderBook.GetBook(types.Sell)

	mustNoErr(t, err)

	restingBuyOrders, exists := buyBook[1000]

	if !exists {
		t.Fatal("expected buy price level 1000 to exist")
	}

	restingSellOrders, exists := sellBook[1100]

	if !exists {
		t.Fatal("expected sell price level 1100 to exist")
	}

	assertRestingOrders(t, restingBuyOrders, 1, 3, 1000, "restingBuyOrders")
	assertRestingOrders(t, restingSellOrders, 1, 1, 1100, "restingSellOrders")

}

func mustNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected error to nil, got %v", err)
	}
}

func assertExecution(t *testing.T, executions []execution.Execution, expectedCount, expectedQty, expectedPrice int) {
	t.Helper()

	if len(executions) != expectedCount {
		t.Fatalf("expected %d execution, got %d", expectedCount, len(executions))
	}

	if len(executions) == 0 {
		return
	}

	if executions[0].Quantity != expectedQty {
		t.Fatalf("expected quantity to be %d, got %d", expectedQty, executions[0].Quantity)
	}

	if executions[0].Price != expectedPrice {
		t.Fatalf("expected price to be %d, got %d", expectedPrice, executions[0].Price)
	}
}

func assertRestingOrders(t *testing.T, orders []types.Order, expectedCount, expectedQty, expectedPrice int, levelName string) {
	t.Helper()

	if len(orders) != expectedCount {
		t.Fatalf("expected %d %s, got %d", expectedCount, levelName, len(orders))
	}

	if len(orders) == 0 {
		return
	}

	if orders[0].Qty != expectedQty {
		t.Fatalf("expected %d quantities, got %d", expectedQty, orders[0].Qty)
	}

	if orders[0].Price != expectedPrice {
		t.Fatalf("expected price to be %d, got %d", expectedPrice, orders[0].Price)
	}
}
