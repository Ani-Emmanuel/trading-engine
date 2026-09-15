package matching

import (
	"cmp"
	"maps"
	"slices"

	orderBook "github.com/ani-emmanuel/trading-engine/orderBook"

	"github.com/ani-emmanuel/trading-engine/execution"
	"github.com/ani-emmanuel/trading-engine/types"
)

type Engine struct {
	OrderBook *orderBook.OrderBook
}

func (e *Engine) Match(takerOrder types.Order) ([]execution.Execution, error) {
	tradeSide := getTradeSide(takerOrder.Side)

	orders, err := e.OrderBook.GetBook(tradeSide)

	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		err := e.OrderBook.AddToOrderBook(
			takerOrder.OwnerId,
			takerOrder.Code,
			takerOrder.Side,
			takerOrder.Qty,
			takerOrder.Price,
		)

		return nil, err
	}

	prices := getPriceLevelsInPriorityOrder(orders, tradeSide)

	var executions []execution.Execution
	var priceExecutions []execution.Execution

	takerOrderQty := takerOrder.Qty

	for _, price := range prices {
		takerOrderQty, priceExecutions, err = e.matchAtPriceLevel(tradeSide, orders[price], takerOrder, takerOrderQty)

		if err != nil {
			return nil, err
		}

		executions = append(executions, priceExecutions...)

		if takerOrderQty == 0 {
			break
		}
	}

	if takerOrderQty > 0 {
		err := e.OrderBook.AddToOrderBook(
			takerOrder.OwnerId,
			takerOrder.Code,
			takerOrder.Side,
			takerOrderQty,
			takerOrder.Price,
		)

		if err != nil {
			return nil, err
		}
	}

	return executions, nil
}

func (e *Engine) matchAtPriceLevel(tradeSide types.Side, makerOrders []types.Order, takerOrder types.Order, takerOrderQty int) (int, []execution.Execution, error) {
	removedCount := 0
	var executions []execution.Execution

	for index, makerOrder := range makerOrders {
		if !canMatch(makerOrder, takerOrder) {
			break
		}

		matchedQty := min(makerOrder.Qty, takerOrderQty)
		makerRemainingOrder := makerOrder.Qty - matchedQty
		takerOrderQty -= matchedQty
		actualIndex := index - removedCount

		err := e.updateOrderBookAfterMatch(
			tradeSide,
			makerOrder,
			makerRemainingOrder,
			actualIndex,
		)

		if err != nil {
			return 0, nil, err
		}

		exec := execution.NewExecution(
			makerOrder.ID,
			takerOrder.ID,
			makerOrder.OwnerId,
			takerOrder.OwnerId,
			takerOrder.Code,
			makerOrder.Price,
			matchedQty,
		)

		executions = append(executions, exec)

		if makerRemainingOrder == 0 {
			removedCount++
		}

		if takerOrderQty == 0 {
			break
		}
	}

	return takerOrderQty, executions, nil
}

func (e *Engine) updateOrderBookAfterMatch(tradeSide types.Side, matchedOrder types.Order, matchedOrderRemaining int, actualIndex int) error {
	if matchedOrderRemaining == 0 {
		return e.OrderBook.RemoveFromOrderBook(tradeSide, matchedOrder.Price, actualIndex)
	}

	return e.OrderBook.UpdateOrderBookRecord(tradeSide, matchedOrderRemaining, matchedOrder.Price, actualIndex)
}

func getPriceLevelsInPriorityOrder(order map[int][]types.Order, side types.Side) []int {
	prices := maps.Keys(order)

	if side == types.Buy {
		return slices.Sorted(prices)
	}

	return slices.SortedFunc(prices, func(a, b int) int {
		return cmp.Compare(b, a)
	})
}

func canMatch(matchedOrder, incomingOrder types.Order) bool {
	if incomingOrder.Side == types.Buy {
		return incomingOrder.Price >= matchedOrder.Price
	}

	return incomingOrder.Price <= matchedOrder.Price
}

func getTradeSide(side types.Side) types.Side {
	switch side {
	case types.Buy:
		return types.Sell
	case types.Sell:
		return types.Buy
	default:
		return ""
	}
}
