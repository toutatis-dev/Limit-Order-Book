package main

import (
	"sync"
	"time"
)

type Side int64

const (
	Buy  Side = 1
	Sell Side = 0
)

type PriceLevel struct {
	Head          *OrderNode
	Tail          *OrderNode
	TotalQuantity uint64
}

type OrderNode struct {
	ID       uint64
	Side     Side
	Price    uint64
	Quantity uint64
	Next     *OrderNode
	Prev     *OrderNode
}

type OrderBook struct {
	BestPrice  uint64
	PriceLevel map[uint64]*PriceLevel //key is price level, value is head of doubly linked list at that pricelevel.
	Orders     map[uint64]*OrderNode  //full map of orders is kept for O(1) cancellation.
	mu         sync.Mutex
}

type Engine struct {
	BuyBook  *OrderBook
	SellBook *OrderBook
}

type Trade struct {
	TakerID   uint64
	MakerID   uint64
	Price     uint64
	Quantity  uint64
	Timestamp time.Time
}

func main() {

}

func NewOrder(id uint64, side Side, price uint64, quantity uint64) *OrderNode {

	return &OrderNode{
		ID:       id,
		Side:     side,
		Price:    price,
		Quantity: quantity,
		Next:     nil,
		Prev:     nil,
	}

}

func (ob *OrderBook) AddOrder(order *OrderNode) {
	//take in an order and add to correct book. We assume the engine has already decided which side the order is on. Check if the book is empty; if so make Head & Tail point to this order.
	//If not, search through the book, find the price level and then insert to the bottom of the price level.
}

func (pl PriceLevel) RemoveOrder(order *OrderNode) {

}

func (e *Engine) Process(order *OrderNode) []Trade {
	//takes in an order, filters by side, makes any trades it can and if it is partial or unfulfilled, adds it to the correct book.
	//buy matches if order price >= sellbook ask, sell price matches if <= buy bid

	var oppositeBook *OrderBook
	var ownBook *OrderBook

	if order.Side == Buy {
		ownBook = e.BuyBook
		oppositeBook = e.SellBook
	} else {
		ownBook = e.SellBook
		oppositeBook = e.BuyBook
	}

	if order.Side == Buy && order.Price < oppositeBook.BestPrice {
		ownBook.AddOrder(order)
	}
	if order.Side == Sell && order.Price > oppositeBook.BestPrice {
		ownBook.AddOrder(order)
	}

	//if we are here then the order is a valid match, so we need to find the PriceLevel by index of the order.price, and start decrementing quantity.

	priceLevel := oppositeBook.PriceLevel[oppositeBook.BestPrice]

	current := priceLevel.Head
	trades := []Trade{}
	for current != nil && order.Quantity > 0 {
		var matchQTY uint64
		if current.Quantity >= order.Quantity {
			matchQTY = order.Quantity
		} else {
			matchQTY = current.Quantity
		}

		order.Quantity -= matchQTY
		current.Quantity -= matchQTY

		next := current.Next
		if current.Quantity == 0 {
			priceLevel.RemoveOrder(current)
		}

		trades = append(trades, Trade{order.ID, current.ID, order.Price, matchQTY, time.Now()})

		current = next
	}

}
