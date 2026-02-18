package main

import (
	"fmt"
	"sync"
)

type Side int64

const (
	Buy Side = 1
	Sell Side = 0
)

type OrderNode struct {
	ID	uint64
	Side	Side
	Price	uint64
	Quantity	uint64
	Next 	*OrderNode
	Prev	*OrderNode
}

type OrderBook struct {
	Head	*OrderNode
	Tail	*OrderNode
	Orders	map[uint64]*OrderNode
	mu	sync.Mutex
}

type Engine struct {
	BuyBook	*OrderBook
	SellBook *OrderBook
}

type Trade struct {

}

func main(){

}

func NewOrder(id uint64, side Side, price uint64, quantity uint64) *OrderNode{

	return &OrderNode{
		ID: id,
		Side: side,
		Price: price,
		Quantity: quantity,
		Next: nil,
		Prev: nil,
	}

}