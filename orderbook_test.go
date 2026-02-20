package main

import (
	"testing"
)

func TestNewOrder(t *testing.T) {
	testOrder := NewOrder(1, Buy, 100, 10)
	example := OrderNode{ID: 1, Side: Buy, Price: 100, Quantity: 10}

	if testOrder.Next != example.Next {
		t.Errorf("Expected %v, got %v", example.Next, testOrder.Next)
	}
	if testOrder.Prev != example.Prev {
		t.Errorf("Expected %v, got %v", example.Prev, testOrder.Prev)
	}
	if testOrder.ID != example.ID {
		t.Errorf("Expected %v, got %v", example.ID, testOrder.ID)
	}
	if testOrder.Side != example.Side {
		t.Errorf("Expected %v, got %v", example.Side, testOrder.Side)
	}
	if testOrder.Price != example.Price {
		t.Errorf("Expected %v, got %v", example.Price, testOrder.Price)
	}
	if testOrder.Quantity != example.Quantity {
		t.Errorf("Expected %v, got %v", example.Quantity, testOrder.Quantity)
	}
}

func TestInsertPriceLevel(t *testing.T)  {
	ob := OrderBook{
		PriceLevel: make(map[uint64]*PriceLevel),
		Orders: make(map[uint64]*OrderNode),
		}
		ob.InsertPriceLevel(10)
		ob.InsertPriceLevel(20)

		if ob.SortedPL[0] != 20 {
			t.Errorf("Expected 20, got %v", ob.SortedPL[0])
		}
		if ob.SortedPL[1] != 10 {
			t.Errorf("Expected 20, got %v", ob.SortedPL[1])
		}

		if ob.PriceLevel[10].Head != nil {
			t.Errorf("Expected nil, got %v", ob.PriceLevel[10].Head)
		}
		if ob.PriceLevel[10].Tail != nil {
			t.Errorf("Expected nil, got %v", ob.PriceLevel[10].Tail)
		}
		if ob.PriceLevel[10].TotalQuantity != 0 {
			t.Errorf("Expected 0, got %v", ob.PriceLevel[10].TotalQuantity)
		}
		

		if ob.PriceLevel[20].Head != nil {
			t.Errorf("Expected nil, got %v", ob.PriceLevel[20].Head)
		}
		if ob.PriceLevel[20].Tail != nil {
			t.Errorf("Expected nil, got %v", ob.PriceLevel[20].Tail)
		}
		if ob.PriceLevel[20].TotalQuantity != 0 {
			t.Errorf("Expected 0, got %v", ob.PriceLevel[20].TotalQuantity)
		}
}