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
