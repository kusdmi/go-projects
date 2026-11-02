package ledger

import "fmt"

var transactions []Transaction

func init() {
	transactions = make([]Transaction, 0)
}

func categoryTotal(category string) float64 {
	var sum float64
	for _, tx := range transactions {
		if tx.Category == category {
			sum += tx.Amount
		}
	}
	return sum
}

func AddTransaction(tx Transaction) error {
	// Валидация вынесена в метод Validate()
	if err := tx.Validate(); err != nil {
		return err
	}

	// Проверка бюджета — бизнес-логика Ledger
	if b, ok := GetBudget(tx.Category); ok {
		newTotal := categoryTotal(tx.Category) + tx.Amount
		if newTotal > b.Limit {
			return fmt.Errorf(
				"budget exceeded for category %q: limit=%.2f, would be=%.2f",
				b.Category, b.Limit, newTotal,
			)
		}
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

func Reset() {
	transactions = make([]Transaction, 0)
	budgets = make(map[string]Budget)
}