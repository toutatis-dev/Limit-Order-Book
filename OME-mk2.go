package main

import (
	"sort"
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
	SortedPL   []uint64
	Side       Side
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

func (ob *OrderBook) NextBestPriceLevel() *PriceLevel {
	//filter logic based on side, if bid then highest price is best, if ask - lowest price wins
	//Based on side, check either tail or head of map[]uint64 for SortedPL (organised largest to smallest)
	if ob.SortedPL == nil {
		return nil
	}

	if ob.Side == Buy {
		ob.BestPrice = ob.SortedPL[0]
	} else {
		ob.BestPrice = ob.SortedPL[len(ob.SortedPL)-1]
	}

	return ob.PriceLevel[ob.BestPrice]
}

func (ob *OrderBook) InsertPriceLevel(price uint64) {
	//use sort.Search to find insertion point assuming descending order.
	//sort.Search handles the empty case by returning index 0
	// sort.Search(len(sortedPL), func (i int) bool {return sortedPL[i] <= newPL})

	i := sort.Search(len(ob.SortedPL), func(i int) bool { return ob.SortedPL[i] <= price })

	valAbove := ob.SortedPL[i:]
	valBelow := ob.SortedPL[:i:i]
	newSlice := append(valBelow, price)
	newSlice = append(newSlice, valAbove...)
	ob.SortedPL = newSlice
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

		if next == nil {
			//call NextBestPriceLevel it should return a ptr to the current BestPL so that we can set next to point at the head of the PL
			oppositeBook.removePriceLevel(current.Price)
			nextPL := oppositeBook.NextBestPriceLevel()
			next = nextPL.Head
		}

		current = next

	}
	return trades
}

func (ob *OrderBook) removePriceLevel(price uint64) {
	//delete the price level from the map
	//clear the price level from the sorted list
	delete(ob.PriceLevel, price)
	i := sort.Search(len(ob.SortedPL), func(i int) bool { return ob.SortedPL[i] <= price })

	belowPL := ob.SortedPL[:i]
	abovePL := ob.SortedPL[i+1:]
	ob.SortedPL = append(belowPL, abovePL...)

}
