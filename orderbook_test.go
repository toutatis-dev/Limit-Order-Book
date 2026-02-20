package main

import (
	"testing"
)

func TestNewOrder(t *testing.T) {
	testOrder := NewOrder(1, Buy, 100, 10)
	example := OrderNode{ID: 1, Side: Buy, Price: 100, Quantity: 10}

	if testOrder.Next != example.Next {
		t.Errorf("Next incorrectly init. Expected %v, got %v", example.Next, testOrder.Next)
	}
	if testOrder.Prev != example.Prev {
		t.Errorf("Prev incorrectly init. Expected %v, got %v", example.Prev, testOrder.Prev)
	}
	if testOrder.ID != example.ID {
		t.Errorf("ID incorrectly init. Expected %v, got %v", example.ID, testOrder.ID)
	}
	if testOrder.Side != example.Side {
		t.Errorf("Side incorrectly init. Expected %v, got %v", example.Side, testOrder.Side)
	}
	if testOrder.Price != example.Price {
		t.Errorf("Price incorrectly init. Expected %v, got %v", example.Price, testOrder.Price)
	}
	if testOrder.Quantity != example.Quantity {
		t.Errorf("Quantity incorrectly init. Expected %v, got %v", example.Quantity, testOrder.Quantity)
	}
}

func TestInsertPriceLevel(t *testing.T) {
	ob := NewOrderBook(Buy)
	ob.InsertPriceLevel(10)
	ob.InsertPriceLevel(20)

	if ob.SortedPL[0] != 20 {
		t.Errorf("Sorted PL incorrectly. Expected 20, got %v\n", ob.SortedPL[0])
	}
	if ob.SortedPL[1] != 10 {
		t.Errorf("Sorted PL incorrectly. Expected 20, got %v\n", ob.SortedPL[1])
	}

	if ob.PriceLevel[10].Head != nil {
		t.Errorf("Incorrectly init price level head. Expected nil, got %v\n", ob.PriceLevel[10].Head)
	}
	if ob.PriceLevel[10].Tail != nil {
		t.Errorf("Incorrectly init price level tail. Expected nil, got %v\n", ob.PriceLevel[10].Tail)
	}
	if ob.PriceLevel[10].TotalQuantity != 0 {
		t.Errorf("Incorrectly init price level total quantity. Expected 0, got %v\n", ob.PriceLevel[10].TotalQuantity)
	}

	if ob.PriceLevel[20].Head != nil {
		t.Errorf("Incorrectly init price level head. Expected nil, got %v\n", ob.PriceLevel[20].Head)
	}
	if ob.PriceLevel[20].Tail != nil {
		t.Errorf("Incorrectly init price level tail. Expected nil, got %v\n", ob.PriceLevel[20].Tail)
	}
	if ob.PriceLevel[20].TotalQuantity != 0 {
		t.Errorf("Incorrectly init price level total quantity. Expected 0, got %v\n", ob.PriceLevel[20].TotalQuantity)
	}
}

func TestAddOrder(t *testing.T) {

	ob := NewOrderBook(Buy)

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
	if ob.PriceLevel[testOrder2.Price].Head.Next != testOrder2 { //verify that the second order is second in the price level
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

func TestProcess(t *testing.T) {
	//Initialise an engine, order books and orders.
	//check that first order makes it onto the book
	//construct a complementary order that will go on to the second book
	//third order that will execute a trade fully, fourth order that will partially fill a trade
	//ensure that it moves down/up the price levels correctly

	e := Engine{
		SellBook: NewOrderBook(Sell),
		BuyBook:  NewOrderBook(Buy),
	}
	//bid order 1 - price 100, quant 20 - straight on the books with no match.
	//ask order 1 - price 120, quant 15, - straight on the books with no match.
	//bid order 2 - price 120, quant 15, - full match, never hits the book and empties ask order1
	//ask order 2 - price 100, quant 10, - partial match, will never hit the books itself but will decrement bid order 1
	//ask order 3 - price 100, quant 15, - another partial match but will empty the bid on books, and get sent to books itself
	//ask order 4 - price 110, quant 5, - add another price level to sell book, will be executed after the 100 level
	//bid order 3 - price 120, quant 10 - will completely empty both books, and fall through 2 price levels in the sell book

	bidOrder1 := NewOrder(1, Buy, 100, 20)
	trade := e.Process(bidOrder1)
	//check it hits the book
	if e.BuyBook.PriceLevel[100].Head != bidOrder1 {
		t.Errorf("bidOrder1 did not hit buybook correctly. Expected %v, got %v\n", bidOrder1, e.BuyBook.PriceLevel[100].Head)
	}
	if len(trade) != 0 {
		t.Errorf("A trade was returned after bidOrder1")
	}

	askOrder1 := NewOrder(2, Sell, 120, 15)
	trade = e.Process(askOrder1)
	//check if it hits the book
	if e.SellBook.PriceLevel[120].Head != askOrder1 {
		t.Errorf("askOrder1 did not hit sellbook correctly. Expected %v, got %v\n", askOrder1, e.SellBook.PriceLevel[120].Head)
	}
	if len(trade) != 0 {
		t.Errorf("A trade was returned after askOrder1")
	}

	bidOrder2 := NewOrder(3, Buy, 120, 15)
	expectedQty := bidOrder2.Quantity
	trade = e.Process(bidOrder2)
	//check it empties the sell book and returns a slice of trades with correct ids, prices, and quantities
	if e.BuyBook.PriceLevel[120] != nil {
		t.Errorf("bidOrder2 has hit the books\n")
	}
	if len(trade) != 1 {
		t.Errorf("No trade was returned after bidOrder2")
	}
	if trade[0].TakerID != bidOrder2.ID {
		t.Errorf("bidOrder2 trade returned taker ID %d. Expected %d", trade[0].TakerID, bidOrder2.ID)
	}
	if trade[0].MakerID != askOrder1.ID {
		t.Errorf("bidOrder2 trade returned Maker ID %d. Expected %d", trade[0].MakerID, askOrder1.ID)
	}
	if trade[0].Price != askOrder1.Price {
		t.Errorf("bidOrder2 trade returned Price %d. Expected %d", trade[0].Price, askOrder1.Price)
	}
	if trade[0].Quantity != expectedQty {
		t.Errorf("bidOrder2 trade returned Quantity %d. Expected %d", trade[0].Quantity, expectedQty)
	}
	if trade[0].Timestamp.IsZero() {
		t.Errorf("Time stamp not set on bidOrder2 trade")
	}

	askOrder2 := NewOrder(4, Sell, 100, 10)
	expectedQty = askOrder2.Quantity
	trade = e.Process(askOrder2)
	//check that trades slice is returned, check the maker id matches bid order 1. bid order 1 should be decremented by the correct amount
	if e.SellBook.PriceLevel[100] != nil {
		t.Errorf("askOrder2 has hit the books\n")
	}
	if len(trade) != 1 {
		t.Errorf("No trade was returned after askOrder2\n")
	}
	if trade[0].TakerID != askOrder2.ID {
		t.Errorf("askOrder2 trade returned taker ID %d. Expected %d\n", trade[0].TakerID, askOrder2.ID)
	}
	if trade[0].MakerID != bidOrder1.ID {
		t.Errorf("askOrder2 trade returned Maker ID %d. Expected %d\n", trade[0].MakerID, bidOrder1.ID)
	}
	if trade[0].Price != bidOrder1.Price {
		t.Errorf("askOrder2 trade returned Price %d. Expected %d\n", trade[0].Price, bidOrder1.Price)
	}
	if trade[0].Quantity != expectedQty {
		t.Errorf("bidOrder2 trade returned Quantity %d. Expected %d\n", trade[0].Quantity, expectedQty)
	}
	if trade[0].Timestamp.IsZero() {
		t.Errorf("Time stamp not set on bidOrder2 trade\n")
	}

	askOrder3 := NewOrder(5, Sell, 100, 15)
	trade = e.Process(askOrder3)
	//check that buy book is exhausted, and that askorder3 is written to sell book.
	if len(e.BuyBook.SortedPL) != 0 {
		t.Errorf("bidOrder1 was not removed from buy book.\n")
	}
	if len(e.SellBook.SortedPL) == 0 {
		t.Errorf("askOrder3 was not written to sell Book.\n")
	}
	if e.BuyBook.PriceLevel[bidOrder1.Price] != nil {
		t.Errorf("Price level still exists in buy book\n")
	}
	if len(trade) == 0 {
		t.Errorf("no trade executed between askOrder3 and bidOrder1\n")
	}
	if trade[0].Quantity != 10 {
		t.Errorf("Trade not executed properly between askOrder3 and bidOrder1, expected quantity 10, got %d\n", trade[0].Quantity)
	}
	if e.SellBook.PriceLevel[askOrder3.Price].Head.Quantity != 5 {
		t.Errorf("ask Order 3 was not decremented properly. expected 5, got %d\n", e.SellBook.PriceLevel[askOrder3.Price].Head.Quantity)
	}

	askOrder4 := NewOrder(6, Sell, 110, 5)
	trade = e.Process(askOrder4)
	//ask order 4 - price 110, quant 5, - add another price level to sell book, will be executed after the 100 level
	if len(e.SellBook.SortedPL) != 2 {
		t.Errorf("failed to add a new price level for askOrder4\n")
	}
	if e.SellBook.SortedPL[0] != 110 {
		t.Errorf("askOrder4 price added in wrong order\n")
	}
	if len(trade) != 0 {
		t.Errorf("Trade was carried out prematurely for askOrder4\n")
	}

	bidOrder3 := NewOrder(7, Buy, 120, 10)
	expectedQtyAsk3 := askOrder3.Quantity
	expectedQtyAsk4 := askOrder4.Quantity
	trade = e.Process(bidOrder3)
	//bid order 3 - price 120, quant 10 - will completely empty both books, and fall through 2 price levels in the sell book
	//will execute askOrders 3 & 4
	// is it making the trade for full quant of askOrder3?
	if trade[0].MakerID != askOrder3.ID {
		t.Errorf("Trade not executed between bid3 and ask3. Expected %d, got %d\n", askOrder3.ID, trade[0].MakerID)
	}
	if trade[0].Quantity != expectedQtyAsk3 {
		t.Errorf("Wrong quantity for trade between bid3 and ask3. Expected %d, got %d\n", expectedQtyAsk3, trade[0].Quantity)
	}
	//is it making the trade for full quant of askOrder4?
	if len(trade) >= 2 {
		if trade[1].MakerID != askOrder4.ID {
			t.Errorf("Trade not executed between bid3 and ask4. Expected %d, got %d\n", askOrder4.ID, trade[1].MakerID)
		}
		if trade[1].Quantity != expectedQtyAsk4 {
			t.Errorf("Wrong quantity for trade between bid3 and ask4. Expected %d, got %d\n", expectedQtyAsk4, trade[1].Quantity)
		}
	}
	//are both sell side price levels removed?
	if len(e.SellBook.SortedPL) != 0 {
		t.Errorf("Price levels on sell side have not been removed after bid3 trades.")
	}
	//is the bidOrder fully emptied and not written to book
	if len(e.BuyBook.SortedPL) != 0 {
		t.Errorf("Buy side still has price levels after bid3 trades.")
	}

}
