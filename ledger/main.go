package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Ledger service started")

	samples := []Transaction{
		{Amount: 1500.50, Category: "food", Description: "Groceries", Date: "2026-10-05"},
		{Amount: 300, Category: "transport", Description: "Taxi", Date: "2026-10-05"},
		{Amount: 75.20, Category: "entertainment", Description: "Cinema", Date: "2026-10-04"},
	}

	for _, tx := range samples {
		if err := AddTransaction(tx); err != nil {
			log.Printf("failed to add transaction: %v", err)
		}
	}

	if err := AddTransaction(Transaction{Amount: 0, Category: "invalid"}); err != nil {
		fmt.Printf("AddTransaction error: %v\n", err)
	}

	fmt.Println("Transactions:")
	for _, tx := range ListTransactions() {
		fmt.Printf(
			"ID=%d, Amount=%.2f, Category=%s, Description=%s, Date=%s\n",
			tx.ID, tx.Amount, tx.Category, tx.Description, tx.Date,
		)
	}
}
