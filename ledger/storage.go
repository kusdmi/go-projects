package ledger

import (
	"database/sql"
	"fmt"

	"example.com/finance/ledger/internal/db"
)

func categorySpent(category string) (float64, error) {
	var spent float64
	query := `SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE category = $1`
	err := db.DB.QueryRow(query, category).Scan(&spent)
	return spent, err
}

// AddTransaction принимает указатель, чтобы сохранить сгенерированный в БД ID
func AddTransaction(tx *Transaction) error {
	if err := tx.Validate(); err != nil {
		return err
	}

	// Проверяем, есть ли бюджет для этой категории
	budget, err := GetBudget(tx.Category)
	if err == nil {
		spent, err := categorySpent(tx.Category)
		if err != nil {
			return fmt.Errorf("failed to calculate category spending: %w", err)
		}
		if spent+tx.Amount > budget.Limit {
			return fmt.Errorf(
				"budget exceeded for category %q: limit=%.2f, would be=%.2f",
				budget.Category, budget.Limit, spent+tx.Amount,
			)
		}
	} else if err != ErrBudgetNotFound {
		return fmt.Errorf("failed to check budget: %w", err)
	}

	// Вставляем транзакцию и получаем ID
	query := `
		INSERT INTO expenses (amount, category, description, date)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err = db.DB.QueryRow(query, tx.Amount, tx.Category, tx.Description, tx.Date).Scan(&tx.ID)
	return err
}

func ListTransactions() ([]Transaction, error) {
	rows, err := db.DB.Query(`SELECT id, amount, category, description, date FROM expenses ORDER BY date DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []Transaction
	for rows.Next() {
		var tx Transaction
		var dateStr string
		var desc sql.NullString
		if err := rows.Scan(&tx.ID, &tx.Amount, &tx.Category, &desc, &dateStr); err != nil {
			return nil, err
		}
		tx.Description = desc.String
		tx.Date = dateStr
		txs = append(txs, tx)
	}
	return txs, rows.Err()
}

// Reset очищает таблицы в БД (полезно для тестов)
func Reset() error {
	// Очищаем таблицу расходов и сбрасываем счетчик ID
	_, err1 := db.DB.Exec("TRUNCATE TABLE expenses RESTART IDENTITY CASCADE")
	// Очищаем таблицу бюджетов и сбрасываем счетчик ID
	_, err2 := db.DB.Exec("TRUNCATE TABLE budgets RESTART IDENTITY CASCADE")
	
	if err1 != nil {
		return err1
	}
	return err2
}