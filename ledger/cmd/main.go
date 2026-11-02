package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"example.com/finance/ledger"
)

func main() {
	fmt.Println("Ledger service started")

	// 1. Загрузка бюджетов из файла
	file, err := os.Open("budgets.json")
	if err != nil {
		log.Fatalf("failed to open budgets.json: %v", err)
	}
	defer file.Close()

	if err := ledger.LoadBudgets(bufio.NewReader(file)); err != nil {
		log.Fatalf("failed to load budgets: %v", err)
	}
	fmt.Println("Budgets loaded from budgets.json")

	// 2. Установка дополнительного бюджета через SetBudget
	if err := ledger.SetBudget(ledger.Budget{Category: "transport", Limit: 1000, Period: "month"}); err != nil {
		log.Printf("failed to set budget: %v", err)
	}
	fmt.Println("Budget for transport set via SetBudget")

	// 3. Демонстрация полиморфизма через интерфейс Validatable
	fmt.Println("--- Validatable interface demonstration ---")
	validTx := ledger.Transaction{Amount: 100, Category: "food"}
	invalidTx := ledger.Transaction{Amount: -10, Category: "food"}
	validBudget := ledger.Budget{Category: "food", Limit: 1000}
	invalidBudget := ledger.Budget{Category: "", Limit: -50}

	fmt.Printf("CheckValid(validTx)      err=%v\n", ledger.CheckValid(validTx))
	fmt.Printf("CheckValid(invalidTx)    err=%v\n", ledger.CheckValid(invalidTx))
	fmt.Printf("CheckValid(validBudget)  err=%v\n", ledger.CheckValid(validBudget))
	fmt.Printf("CheckValid(invalidBudget) err=%v\n", ledger.CheckValid(invalidBudget))

	// 4. Обычные транзакции
	fmt.Println("--- Adding transactions ---")
	samples := []ledger.Transaction{
		{Amount: 1500.50, Category: "food", Description: "Groceries", Date: "2026-10-05"},
		{Amount: 300, Category: "transport", Description: "Taxi", Date: "2026-10-05"},
		{Amount: 75.20, Category: "entertainment", Description: "Cinema", Date: "2026-10-04"},
	}
	for _, tx := range samples {
		if err := ledger.AddTransaction(tx); err != nil {
			fmt.Printf("AddTransaction rejected: %v\n", err)
			continue
		}
		fmt.Printf("AddTransaction accepted: category=%s amount=%.2f\n", tx.Category, tx.Amount)
	}

	// 5. Транзакция сверх лимита бюджета food (5000)
	overBudget := ledger.Transaction{
		Amount: 4000, Category: "food", Description: "Big shopping", Date: "2026-10-05",
	}
	if err := ledger.AddTransaction(overBudget); err != nil {
		fmt.Printf("AddTransaction rejected (budget exceeded expected): %v\n", err)
	}

	// 6. Некорректная транзакция (отрицательная сумма)
	if err := ledger.AddTransaction(ledger.Transaction{Amount: -1, Category: "food"}); err != nil {
		fmt.Printf("AddTransaction rejected (negative amount expected): %v\n", err)
	}

	// 7. Вывод списка транзакций
	fmt.Println("--- Transactions ---")
	for _, tx := range ledger.ListTransactions() {
		fmt.Printf(
			"ID=%d, Amount=%.2f, Category=%s, Description=%s, Date=%s\n",
			tx.ID, tx.Amount, tx.Category, tx.Description, tx.Date,
		)
	}
}
