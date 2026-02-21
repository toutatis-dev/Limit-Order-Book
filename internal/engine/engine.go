package engine

import (
	"fmt"
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
	TakerID   uint64    `json:"taker_id"`
	MakerID   uint64    `json:"maker_id"`
	Price     uint64    `json:"price"`
	Quantity  uint64    `json:"quantity"`
	Timestamp time.Time `json:"timestamp"`
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

func (ob *OrderBook) AddOrder(order *OrderNode) error {
	//take the order and discover the price level. search the map in o(1) for the price level. if it doesnt exist, then create the price level and point the head & tail to this order.
	//make sure to add the new price level to the sorted slice.
	//if the price level exists, the current tail next needs to point to this order, the order.Prev needs to point to the tail and order.Next to nil, then update price level tail to order.
	//any order must be added to the order map for cancellation and best price updated - remember to also ensure that bestAsk is != 0.

	_, exists := ob.Orders[order.ID]
	if exists {
		return fmt.Errorf("Duplicate Order ID: %d", order.ID)
	}

	pl := ob.PriceLevel[order.Price]

	if pl == nil {
		ob.InsertPriceLevel(order.Price)
		pl = ob.PriceLevel[order.Price]
		pl.Head = order
		pl.Tail = order
	} else {
		order.Prev = pl.Tail
		order.Next = nil
		pl.Tail.Next = order
		pl.Tail = order
	}

	pl.TotalQuantity += order.Quantity

	ob.Orders[order.ID] = order
	if len(ob.SortedPL) == 1 {
		ob.BestPrice = order.Price
	} else {
		if ob.Side == Buy && order.Price > ob.BestPrice {
			ob.BestPrice = order.Price
		}
		if ob.Side == Sell && (order.Price < ob.BestPrice || ob.BestPrice == 0) {
			ob.BestPrice = order.Price
		}
	}
	return nil
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

	if len(ob.SortedPL) == 0 {
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

	ob.PriceLevel[price] = &PriceLevel{Head: nil, Tail: nil, TotalQuantity: 0}
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

	if len(oppositeBook.SortedPL) == 0 {
		//no orders on opposite book.
		ownBook.AddOrder(order)
		return []Trade{}
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
		priceLevel.TotalQuantity -= matchQTY

		trades = append(trades, Trade{order.ID, current.ID, current.Price, matchQTY, time.Now()})

		next := current.Next
		if current.Quantity == 0 {
			removedOrder := priceLevel.RemoveHead()
			if removedOrder != nil {
				delete(oppositeBook.Orders, removedOrder.ID)
			}

		}

		if next == nil && priceLevel.TotalQuantity == 0 {
			//call NextBestPriceLevel it should return a ptr to the current BestPL so that we can set next to point at the head of the PL

			oppositeBook.removePriceLevel(current.Price)
			nextPL := oppositeBook.NextBestPriceLevel()
			if nextPL == nil {

				if order.Quantity > 0 {
					ownBook.AddOrder(order)
				}
				return trades
			}

			if order.Side == Buy && order.Price < oppositeBook.BestPrice {
				if order.Quantity > 0 {
					ownBook.AddOrder(order)
				}
				return trades
			}
			if order.Side == Sell && order.Price > oppositeBook.BestPrice {
				if order.Quantity > 0 {
					ownBook.AddOrder(order)
				}
				return trades
			}
			next = nextPL.Head
			priceLevel = nextPL
		}

		current = next

	}

	if order.Quantity > 0 {
		ownBook.AddOrder(order)
	}
	return trades
}

func (ob *OrderBook) removePriceLevel(price uint64) {
	//delete the price level from the map
	//clear the price level from the sorted list
	delete(ob.PriceLevel, price)
	i := sort.Search(len(ob.SortedPL), func(i int) bool { return ob.SortedPL[i] <= price })

	if i < len(ob.SortedPL) && ob.SortedPL[i] == price {
		ob.SortedPL = append(ob.SortedPL[:i], ob.SortedPL[i+1:]...)
	}

}

func NewOrderBook(side Side) *OrderBook {
	return &OrderBook{
		PriceLevel: make(map[uint64]*PriceLevel),
		Orders:     make(map[uint64]*OrderNode),
		Side:       side,
	}
}

func (e *Engine) CancelOrder(id uint64) error {

	if e.BuyBook.Orders[id] != nil {
		return e.BuyBook.CancelOrder(id)
	}
	if e.SellBook.Orders[id] != nil {
		return e.SellBook.CancelOrder(id)
	}
	return fmt.Errorf("ID %d not in either book", id)

}

func (ob *OrderBook) CancelOrder(id uint64) error {
	//search the id in the Orders map to find the node.
	//set order.Prev.Next to order.Next
	//set Order.Next.Prev to order.Prev
	//we also need to decrement the pl.totalquantity by order.quantity
	//if either order.Prev or order.Next are nil, then its at the end or the beginning of the PL.
	//if both are nil, then the price level is empty, so we must pass it to ob.RemovePriceLevel
	//find the price level by looking up order.Price in the pl map. and update pl.Head or pl.Tail accordingly.
	//then we should call ob.NextBestPriceLevel, to make sure that it is updated.

	order, ok := ob.Orders[id]
	if !ok {
		return fmt.Errorf("ID %d does not exist in book", id)
	}

	//Check if price level will be empty after removal of order.If so, remove it
	if order.Prev == nil && order.Next == nil {
		ob.removePriceLevel(order.Price)
		if order.Price == ob.BestPrice { //if this price level was the best, then determine nextbest
			ob.NextBestPriceLevel()
		}
		delete(ob.Orders, id)
		return nil
	}

	//if order.next is nil, then its at the tail, so pl.tail must point to order.prev
	if order.Next == nil {
		ob.PriceLevel[order.Price].Tail = order.Prev
	} else { //its not at the tail
		order.Next.Prev = order.Prev
	}

	//if order.prev is nil, then its at the head, so pl.head must point to order.next
	if order.Prev == nil {
		ob.PriceLevel[order.Price].Head = order.Next
	} else { //its not at the head
		order.Prev.Next = order.Next
	}

	order.Next = nil
	order.Prev = nil

	ob.PriceLevel[order.Price].TotalQuantity -= order.Quantity
	delete(ob.Orders, id)
	return nil

}

func NewEngine() *Engine {
	return &Engine{
		SellBook: NewOrderBook(Sell),
		BuyBook:  NewOrderBook(Buy),
	}
}
