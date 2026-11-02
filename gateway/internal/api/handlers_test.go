package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/finance/ledger"
)

func TestBudgetAndTransactionFlow(t *testing.T) {
	router := NewRouter()

	// --- Тесты для бюджетов ---
	t.Run("create budget", func(t *testing.T) {
		ledger.Reset()

		body := `{"category":"food","limit":5000}`
		req := httptest.NewRequest(http.MethodPost, "/api/budgets", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}

		contentType := w.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			t.Errorf("expected application/json, got %s", contentType)
		}

		var resp BudgetResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Category != "food" || resp.Limit != 5000 {
			t.Errorf("unexpected response body: %+v", resp)
		}
	})

	t.Run("get budgets", func(t *testing.T) {
		ledger.Reset()
		_ = ledger.SetBudget(ledger.Budget{Category: "food", Limit: 5000})

		req := httptest.NewRequest(http.MethodGet, "/api/budgets", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp []BudgetResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(resp) != 1 || resp[0].Category != "food" || resp[0].Limit != 5000 {
			t.Errorf("unexpected response body: %+v", resp)
		}
	})

	t.Run("create budget bad request", func(t *testing.T) {
		ledger.Reset()

		body := `{"category":"","limit":-100}` // Невалидные данные
		req := httptest.NewRequest(http.MethodPost, "/api/budgets", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	// --- Тесты для транзакций ---
	t.Run("create valid transaction", func(t *testing.T) {
		ledger.Reset()
		_ = ledger.SetBudget(ledger.Budget{Category: "food", Limit: 5000})

		body := `{"amount":450,"category":"food","description":"lunch","date":"2025-09-10"}`
		req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}

		var resp TransactionResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if resp.Amount != 450 || resp.Category != "food" {
			t.Errorf("unexpected response body: %+v", resp)
		}
	})

	t.Run("get transactions", func(t *testing.T) {
		ledger.Reset()
		_ = ledger.AddTransaction(ledger.Transaction{Amount: 450, Category: "food", Date: "2025-09-10"})

		req := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp []TransactionResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(resp) != 1 || resp[0].Amount != 450 || resp[0].Category != "food" {
			t.Errorf("unexpected response body: %+v", resp)
		}
	})

	t.Run("budget exceeded", func(t *testing.T) {
		ledger.Reset()
		_ = ledger.SetBudget(ledger.Budget{Category: "food", Limit: 5000})
		_ = ledger.AddTransaction(ledger.Transaction{Amount: 450, Category: "food"})

		body := `{"amount":5000,"category":"food","description":"shopping","date":"2025-09-11"}`
		req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("expected status 409, got %d", w.Code)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		ledger.Reset()

		body := `{"amount": -50, "category": "food"` // Синтаксическая ошибка JSON
		req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}