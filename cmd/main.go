package main

import (
	"lob/internal/engine"
)

func main() {
	order1 := engine.NewOrder(1, engine.Buy, 100, 10)
	e := engine.Engine{
		SellBook: engine.NewOrderBook(engine.Sell),
		BuyBook:  engine.NewOrderBook(engine.Buy),
	}

	e.Process(order1)
}
