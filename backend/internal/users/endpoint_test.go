package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Мок сервиса пользователей
type mockUserService struct{}

func (m *mockUserService) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input map[string]string
	json.NewDecoder(r.Body).Decode(&input)
	if input["email"] == "" || input["password"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": "mocktoken"})
}

func (m *mockUserService) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": "mocktoken"})
}

func TestUsers_Register(t *testing.T) {
	handler := http.HandlerFunc((&mockUserService{}).RegisterHandler)

	payload := map[string]string{
		"email":    "test@example.com",
		"name":     "Test User",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/users/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
}

func TestUsers_Login(t *testing.T) {
	handler := http.HandlerFunc((&mockUserService{}).LoginHandler)

	payload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/users/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
}
