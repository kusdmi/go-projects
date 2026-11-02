package ledger

import (
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Инициализируем БД для тестов
	if err := Init(); err != nil {
		panic("Failed to init DB for tests: " + err.Error())
	}
	os.Exit(m.Run())
}

func TestTransactionValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tx      Transaction
		wantErr bool
	}{
		{"valid", Transaction{Amount: 100, Category: "food"}, false},
		{"zero amount", Transaction{Amount: 0, Category: "food"}, true},
		{"negative amount", Transaction{Amount: -10, Category: "food"}, true},
		{"empty category", Transaction{Amount: 100, Category: ""}, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.tx.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Transaction.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBudgetValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		b       Budget
		wantErr bool
	}{
		{"valid", Budget{Category: "food", Limit: 1000}, false},
		{"zero limit", Budget{Category: "food", Limit: 0}, true},
		{"negative limit", Budget{Category: "food", Limit: -100}, true},
		{"empty category", Budget{Category: "", Limit: 1000}, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.b.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Budget.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddTransactionWithBudget(t *testing.T) {
	_ = Reset() // Очищаем БД

	if err := SetBudget(Budget{Category: "food", Limit: 5000}); err != nil {
		t.Fatalf("failed to set budget: %v", err)
	}

	tx1 := Transaction{Amount: 1500, Category: "food", Date: "2025-09-10"}
	if err := AddTransaction(&tx1); err != nil { // Передаем указатель
		t.Errorf("expected success, got error: %v", err)
	}

	txs, _ := ListTransactions() // Возвращает 2 значения
	if len(txs) != 1 {
		t.Errorf("expected 1 transaction, got %d", len(txs))
	}

	tx2 := Transaction{Amount: 4000, Category: "food", Date: "2025-09-11"}
	err := AddTransaction(&tx2) // Передаем указатель
	if err == nil {
		t.Fatal("expected budget exceeded error, got nil")
	}
	if !strings.Contains(err.Error(), "budget exceeded") {
		t.Errorf("expected 'budget exceeded' error, got: %v", err)
	}

	txs, _ = ListTransactions()
	if len(txs) != 1 {
		t.Errorf("expected 1 transaction after rejection, got %d", len(txs))
	}
}