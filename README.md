# Trading Engine

A trading engine built in Go, featuring an order book, price-time priority matching, partial order fills, and trade execution.

The project explores the core mechanics behind exchange matching engines, focusing on efficient order management and deterministic trade execution.

## Features

- **Order Book Management** — Maintains buy and sell orders organized by price levels.
- **Order Matching** — Matches incoming orders against compatible resting orders.
- **Price-Time Priority** — Prioritizes better prices and executes older orders first when prices are equal.
- **Partial Order Fills** — Supports orders that are partially executed, with remaining quantities retained in the order book.
- **Trade Execution** — Records executions containing order IDs, participants, execution price, and quantity.
- **Unit Testing** — Tests core matching behavior, order book operations, and execution scenarios.

## How It Works

The engine uses two types of orders:

- **Maker order:** An order already resting in the order book.
- **Taker order:** A new incoming order that attempts to match against existing orders.

When a new order arrives, the engine:

1. Identifies the opposite side of the order book.
2. Sorts available price levels by matching priority.
3. Checks whether the incoming order's price can match existing orders.
4. Executes trades at the resting maker order's price.
5. Updates or removes filled maker orders.
6. Adds any unfilled quantity of the incoming order to the order book.

### Example

Suppose the order book contains a sell order for 2 BTC at 100,000.

A buyer submits an order for 3 BTC at 100,000.

The engine executes 2 BTC and places the remaining 1 BTC buy order in the order book, provided the remaining order is valid.

## Project Structure

```text
.
├── execution/    # Trade execution records
│ ├── execution.go
│ └── ...
├── matching/     # Order matching logic
│ ├── engine.go
│ └── engine_test.go
├── orderBook/    # Buy and sell order management
│ ├── order_book.go
│ └── order_book_test.go
└── types/        # Shared domain types
│ └── types.go
├── go.mod
└── go.sum
```

## Requirements

- Go (version compatible with the project's `go.mod`)

## Getting Started

Clone the repository:

```bash
git clone https://github.com/Ani-Emmanuel/trading-engine.git
cd trading-engine
```

Run the tests:

```bash
go test ./...
```

## Design Principles

- Keep order book management separate from matching logic.
- Preserve price-time priority during matching.
- Correctly handle full fills, partial fills, and remaining quantities.
- Keep trade execution records separate from order book state.
- Use automated tests to validate matching behavior.

## Roadmap

- [ ] Complete core order matching and order book functionality.
- [ ] Strengthen validation for invalid orders and unsupported sides.
- [ ] Add tests for price-crossing rules and multiple price levels.
- [ ] Support multiple trading pairs with isolated order books.
- [ ] Add order cancellation and modification.
- [ ] Introduce concurrent order processing and synchronization.
- [ ] Add performance benchmarks.
- [ ] Explore trade settlement and balance management.

## Motivation

This project is an exploration of the engineering principles behind cryptocurrency exchange infrastructure. The goal is to understand how order books, matching algorithms, and trade execution fit together while building a maintainable and testable system in Go.

## Disclaimer

This project is for educational and experimentation purposes. It is not intended for production trading or the management of real funds.
