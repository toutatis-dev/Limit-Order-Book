package api

import (
	"encoding/json"
	"lob/internal/engine"
	"net/http"
	"sync/atomic"
)

type OrderProcessor interface {
	Process(order *engine.OrderNode) []engine.Trade
	CancelOrder(id uint64) error
	GetBookState() engine.BookState
}

type API struct {
	Processor OrderProcessor
	NextID    atomic.Uint64
}

type OrderRequest struct {
	Side     engine.Side `json:"side"`
	Price    uint64      `json:"price"`
	Quantity uint64      `json:"quantity"`
}

type OrderResponse struct {
	ID     uint64         `json:"id"`
	Status string         `json:"status"`
	Trades []engine.Trade `json:"trades"`
}

type CancelRequest struct {
	ID uint64 `json:"id"`
}

type CancelResponse struct {
	ID     uint64 `json:"id"`
	Status string `json:"status"`
}

func NewAPI(e OrderProcessor) *API {
	return &API{
		Processor: e,
	}
}

func (a *API) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderRequest OrderRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&orderRequest)
	if err != nil {
		http.Error(w, "Could not decode JSON", http.StatusBadRequest)
		return
	}
	if orderRequest.Price == 0 {
		http.Error(w, "Price must be greater than 0", http.StatusBadRequest)
		return
	}
	if orderRequest.Quantity == 0 {
		http.Error(w, "Quantity must be greater than 0", http.StatusBadRequest)
		return
	}
	if orderRequest.Side != engine.Buy && orderRequest.Side != engine.Sell {
		http.Error(w, "Side must be either 0 or 1", http.StatusBadRequest)
		return
	}
	//id, side, price, quant
	newOrder := engine.NewOrder(a.NextID.Add(1), orderRequest.Side, orderRequest.Price, orderRequest.Quantity)
	originalQuant := newOrder.Quantity
	trades := a.Processor.Process(newOrder)
	var status string
	if newOrder.Quantity == originalQuant {
		status = "On book"
	} else if newOrder.Quantity < originalQuant && newOrder.Quantity > 0 {
		status = "Partial"
	} else {
		status = "filled"
	}

	response := OrderResponse{
		ID:     newOrder.ID,
		Status: status,
		Trades: trades,
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJson)

}

func (a *API) HandleCancelOrder(w http.ResponseWriter, r *http.Request) {

	var cancelRequest CancelRequest

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&cancelRequest)
	if err != nil {
		http.Error(w, "Could not decode JSON", http.StatusBadRequest)
		return
	}
	err = a.Processor.CancelOrder(cancelRequest.ID)
	if err != nil {
		http.Error(w, "Could not cancel order, check that the ID is valid", http.StatusNotFound)
		return
	}

	cancelResponse := CancelResponse{
		ID:     cancelRequest.ID,
		Status: "Order has been cancelled",
	}

	response, err := json.Marshal(cancelResponse)

	if err != nil {
		http.Error(w, "something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (a *API) HandleGetBook(w http.ResponseWriter, r *http.Request) {
	getBookStateResponse := a.Processor.GetBookState()
	response, err := json.Marshal(getBookStateResponse)
	if err != nil {
		http.Error(w, "Something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}
