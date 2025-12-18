package orders

import (
	"bytes"
	
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Мок сервиса заказов
type mockOrderService struct{}

func (m *mockOrderService) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var input map[string]interface{}
	json.NewDecoder(r.Body).Decode(&input)
	if input["items"] == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "no items"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (m *mockOrderService) ListByUserHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]map[string]interface{}{
		{"id": 1, "user_id": 1, "total_price": 150},
	})
}

func TestOrders_CreateOrder(t *testing.T) {
	handler := http.HandlerFunc((&mockOrderService{}).CreateOrderHandler)

	payload := map[string]interface{}{
		"items": []map[string]interface{}{
			{"product_id": 1, "quantity": 2, "unit_price": 50},
		},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
}

func TestOrders_ListByUser(t *testing.T) {
	handler := http.HandlerFunc((&mockOrderService{}).ListByUserHandler)

	req := httptest.NewRequest("GET", "/orders/me?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
}
