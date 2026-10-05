package main

import "errors"

var transactions []Transaction

func init() {
	transactions = make([]Transaction, 0)
}

func AddTransaction(tx Transaction) error {
	if tx.Amount == 0 {
		return errors.New("amount cannot be zero")
	}

	tx.ID = len(transactions) + 1
	transactions = append(transactions, tx)

	return nil
}

func ListTransactions() []Transaction {
	result := make([]Transaction, len(transactions))
	copy(result, transactions)

	return result
}
