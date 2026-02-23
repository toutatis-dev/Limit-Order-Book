package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"lob/internal/engine"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockEngine struct {
	MockCancelError        error
	MockProcessResult      []engine.Trade
	MockGetBookStateResult engine.BookState
}

func (m *MockEngine) CancelOrder(id uint64) error {

	return m.MockCancelError
}

func (m *MockEngine) Process(order *engine.OrderNode) []engine.Trade {
	return m.MockProcessResult
}

func (m *MockEngine) GetBookState() engine.BookState {
	return m.MockGetBookStateResult
}

func TestHandleCreateOrder(t *testing.T) {

	tests := []struct {
		name           string
		payload        map[string]any
		expectedStatus int
	}{
		{
			name: "Valid Order",
			payload: map[string]any{
				"side":     engine.Buy,
				"price":    10,
				"quantity": 100,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid Order - Zero Quantity",
			payload: map[string]any{
				"side":     engine.Buy,
				"price":    100,
				"quantity": 0,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid Order - Zero Price",
			payload: map[string]any{
				"side":     engine.Sell,
				"price":    0,
				"quantity": 100,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid Order - Invalid Side",
			payload: map[string]any{
				"side":     111,
				"price":    100,
				"quantity": 100,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockEngine{
				MockProcessResult: []engine.Trade{{
					TakerID:   1,
					MakerID:   2,
					Price:     100,
					Quantity:  10,
					Timestamp: time.Now(),
				},
				},
			}
			a := NewAPI(mock)
			body, err := json.Marshal(tt.payload)
			if err != nil {
				t.Errorf("Failed to marshal payload: %v\n", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewBuffer(body))

			rr := httptest.NewRecorder()

			a.HandleCreateOrder(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected %d, got %d", tt.expectedStatus, rr.Code)
			}

		})
	}
}

func TestCancelOrder(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedStatus int
		mockResponse   error
	}{
		{
			name:           "Valid id",
			url:            "/order/1",
			expectedStatus: http.StatusOK,
			mockResponse:   nil,
		},
		{
			name:           "Invalid Cancellation - non numeric id",
			url:            "/order/abc",
			expectedStatus: http.StatusBadRequest,
			mockResponse:   nil,
		},
		{
			name:           "Invalid Cancellation - order ID doesnt exist",
			url:            "/order/2",
			expectedStatus: http.StatusNotFound,
			mockResponse:   fmt.Errorf("Order 2 does could not be found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockEngine{MockCancelError: tt.mockResponse}
			a := NewAPI(mock)
			r := chi.NewRouter()
			r.Delete("/order/{id}", a.HandleCancelOrder)
			req := httptest.NewRequest(http.MethodDelete, tt.url, bytes.NewBuffer(nil))
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
