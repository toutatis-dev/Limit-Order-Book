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
	SortedPl   []*PriceLevel
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

func (pl *PriceLevel) RemoveHead() *OrderNode {
	//removing an order from a price level will always remove from the head, so no need to search list.
	//pL.head.next -> store in var
	// pL.Head -> order .next, unless order.next is nil, in which case set head and tail to nil
	//order.next.prev -> nil
	//clean current order prev/next etc

	if pl.Head == nil {
		return nil
	}

	removedOrder := pl.Head
	if pl.Head.Next == nil {
		pl.Head = nil
		pl.Tail = nil
	} else {
		pl.Head = pl.Head.Next
		pl.Head.Prev = nil
	}
	removedOrder.Next = nil
	removedOrder.Prev = nil

	return removedOrder
}

func (ob *OrderBook) NextBestPriceLevel(side Side) {
	//filter logic based on side, if bid then highest price is best, if ask - lowest price wins
	//loop through the map to find the best price level, update the bestPrice var in the book

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
		return []Trade{}
	}
	if order.Side == Sell && order.Price > oppositeBook.BestPrice {
		ownBook.AddOrder(order)
		return []Trade{}
	}

	//if we are here then the order is a valid match, so we need to find the PriceLevel by index of the order.price, and start decrementing quantity.

	priceLevel := oppositeBook.PriceLevel[oppositeBook.BestPrice]

	current := priceLevel.Head
	trades := []Trade{}
	for current != nil && order.Quantity > 0 { //will exhaust the current price level for as long as there is quantity in the order.
		var matchQTY uint64
		if current.Quantity >= order.Quantity {
			matchQTY = order.Quantity
		} else {
			matchQTY = current.Quantity
		}

		order.Quantity -= matchQTY
		current.Quantity -= matchQTY
		trades = append(trades, Trade{order.ID, current.ID, current.Price, matchQTY, time.Now()})

		next := current.Next
		if current.Quantity == 0 {
			_ = priceLevel.RemoveHead()
		}

		current = next
	}
	return trades
}
