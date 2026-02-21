package api

import (
	"encoding/json"
	"lob/internal/engine"
	"net/http"
	"sync/atomic"
)

type API struct {
	engine *engine.Engine
	nextID atomic.Uint64
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

func NewAPI(e *engine.Engine) *API {
	return &API{
		engine: e,
	}
}

func (a *API) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderRequest OrderRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&orderRequest)
	if err != nil {
		http.Error(w, "something went wrong", 500)
		return
	}
	if orderRequest.Price == 0 {
		http.Error(w, "something went wrong", 500)
		return
	}
	if orderRequest.Quantity == 0 {
		http.Error(w, "something went wrong", 500)
		return
	}
	if orderRequest.Side != engine.Buy && orderRequest.Side != engine.Sell {
		http.Error(w, "something went wrong", 500)
		return
	}
	//id, side, price, quant
	newOrder := engine.NewOrder(a.nextID.Add(1), orderRequest.Side, orderRequest.Price, orderRequest.Quantity)
	originalQuant := newOrder.Quantity
	trades := a.engine.Process(newOrder)
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
