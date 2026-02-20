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
			t.Errorf("Expected 20, got %v\n", ob.SortedPL[0])
		}
		if ob.SortedPL[1] != 10 {
			t.Errorf("Expected 20, got %v\n", ob.SortedPL[1])
		}

		if ob.PriceLevel[10].Head != nil {
			t.Errorf("Expected nil, got %v\n", ob.PriceLevel[10].Head)
		}
		if ob.PriceLevel[10].Tail != nil {
			t.Errorf("Expected nil, got %v\n", ob.PriceLevel[10].Tail)
		}
		if ob.PriceLevel[10].TotalQuantity != 0 {
			t.Errorf("Expected 0, got %v\n", ob.PriceLevel[10].TotalQuantity)
		}
		

		if ob.PriceLevel[20].Head != nil {
			t.Errorf("Expected nil, got %v\n", ob.PriceLevel[20].Head)
		}
		if ob.PriceLevel[20].Tail != nil {
			t.Errorf("Expected nil, got %v\n", ob.PriceLevel[20].Tail)
		}
		if ob.PriceLevel[20].TotalQuantity != 0 {
			t.Errorf("Expected 0, got %v\n", ob.PriceLevel[20].TotalQuantity)
		}
}

func TestAddOrder(t *testing.T)  {
	
	ob := OrderBook{
		PriceLevel: make(map[uint64]*PriceLevel),
		Orders: make(map[uint64]*OrderNode),
		Side: Buy,
		}
	testOrder := NewOrder(1, Buy, 100, 10)

	ob.AddOrder(testOrder)

	//test that the Head and Tail of the created price level point to the new order.
	if ob.PriceLevel[testOrder.Price].Head != testOrder {
		t.Errorf("PL Head not pointing to only order. Expected %v, got %v\n", testOrder, ob.PriceLevel[testOrder.Price].Head)
	}
	if ob.PriceLevel[testOrder.Price].Tail != testOrder {
		t.Errorf("PL Tail not pointing to only order. Expected %v, got %v\n", testOrder, ob.PriceLevel[testOrder.Price].Tail)
	}


	testOrder2 := NewOrder(2, Buy, 100, 20)
	ob.AddOrder(testOrder2)
	//test that the second order is inserted into the price level after the first.
	if ob.PriceLevel[testOrder2.Price].Head.Next != testOrder2 {//verify that the second order is second in the price level
		t.Errorf("Second Order not inserted into price level correctly. Expected %v, got %v\n", testOrder2, ob.PriceLevel[testOrder2.Price].Head.Next)
	}
	//check that thet ob.Orders map is intitialised with the orders for cancellations
	if ob.Orders[testOrder.ID] != testOrder {
		t.Errorf("Orders Map is wrong. Expected %v, got %v\n", testOrder, ob.Orders[testOrder.ID])
	}
	if ob.Orders[testOrder2.ID] != testOrder2 {
		t.Errorf("Orders Map is wrong. Expected %v, got %v\n", testOrder2, ob.Orders[testOrder2.ID])
	}
	//check that the best price is updated.

	if ob.BestPrice != testOrder.Price {
		t.Errorf("BestPrice is wrong. Expected %v, got %v\n", testOrder.Price, ob.BestPrice)
	}

}