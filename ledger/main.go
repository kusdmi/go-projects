package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Ledger service started")

	file, err := os.Open("budgets.json")
	if err != nil {
		log.Fatalf("failed to open budgets.json: %v", err)
	}
	defer file.Close()

	if err := LoadBudgets(bufio.NewReader(file)); err != nil {
		log.Fatalf("failed to load budgets: %v", err)
	}
	fmt.Println("Budgets loaded from budgets.json")

	if err := SetBudget(Budget{Category: "transport", Limit: 1000, Period: "month"}); err != nil {
		log.Printf("failed to set budget: %v", err)
	}
	fmt.Println("Budget for transport set via SetBudget")

	samples := []Transaction{
		{Amount: 1500.50, Category: "food", Description: "Groceries", Date: "2026-10-05"},
		{Amount: 300, Category: "transport", Description: "Taxi", Date: "2026-10-05"},
		{Amount: 75.20, Category: "entertainment", Description: "Cinema", Date: "2026-10-04"},
	}
	for _, tx := range samples {
		if err := AddTransaction(tx); err != nil {
			fmt.Printf("AddTransaction rejected: %v\n", err)
			continue
		}
		fmt.Printf("AddTransaction accepted: category=%s amount=%.2f\n", tx.Category, tx.Amount)
	}

	overBudget := Transaction{
		Amount: 4000, Category: "food", Description: "Big shopping", Date: "2026-10-05",
	}
	if err := AddTransaction(overBudget); err != nil {
		fmt.Printf("AddTransaction rejected (budget exceeded expected): %v\n", err)
	}

	if err := AddTransaction(Transaction{Amount: 0, Category: "food"}); err != nil {
		fmt.Printf("AddTransaction rejected (zero amount expected): %v\n", err)
	}

	fmt.Println("Transactions:")
	for _, tx := range ListTransactions() {
		fmt.Printf(
			"ID=%d, Amount=%.2f, Category=%s, Description=%s, Date=%s\n",
			tx.ID, tx.Amount, tx.Category, tx.Description, tx.Date,
		)
	}
}