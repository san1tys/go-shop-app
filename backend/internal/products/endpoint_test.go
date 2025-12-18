package products

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Мок сервиса продуктов
type mockProductService struct{}

func (m *mockProductService) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]map[string]interface{}{
		{"id": 1, "name": "Product 1", "price": 100},
	})
}

func (m *mockProductService) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": 1, "name": "Product 1", "price": 100})
}

func TestProducts_GetAll(t *testing.T) {
	handler := http.HandlerFunc((&mockProductService{}).GetAllHandler)

	req := httptest.NewRequest("GET", "/products?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
}

func TestProducts_GetByID(t *testing.T) {
	handler := http.HandlerFunc((&mockProductService{}).GetByIDHandler)

	req := httptest.NewRequest("GET", "/products/1", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
}
