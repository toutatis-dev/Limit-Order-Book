package main

import (
	"sync"
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
