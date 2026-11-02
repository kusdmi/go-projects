package ledger

import (
	"database/sql"
	"errors"
	"fmt"

	"example.com/finance/ledger/internal/db"
)

type Budget struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
	Period   string  `json:"period,omitempty"`
}

func (b Budget) Validate() error {
	if b.Category == "" {
		return errors.New("budget category cannot be empty")
	}
	if b.Limit <= 0 {
		return fmt.Errorf("budget limit must be positive, got %.2f", b.Limit)
	}
	return nil
}

var ErrBudgetNotFound = errors.New("budget not found")

func SetBudget(b Budget) error {
	if err := b.Validate(); err != nil {
		return err
	}
	query := `
		INSERT INTO budgets (category, limit_amount)
		VALUES ($1, $2)
		ON CONFLICT (category) DO UPDATE
		SET limit_amount = EXCLUDED.limit_amount
	`
	_, err := db.DB.Exec(query, b.Category, b.Limit)
	return err
}

func GetBudget(category string) (Budget, error) {
	var b Budget
	query := `SELECT category, limit_amount FROM budgets WHERE category = $1`
	err := db.DB.QueryRow(query, category).Scan(&b.Category, &b.Limit)
	if err == sql.ErrNoRows {
		return Budget{}, ErrBudgetNotFound
	}
	return b, err
}

func ListBudgets() ([]Budget, error) {
	rows, err := db.DB.Query(`SELECT category, limit_amount FROM budgets ORDER BY category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []Budget
	for rows.Next() {
		var b Budget
		if err := rows.Scan(&b.Category, &b.Limit); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}
	return budgets, rows.Err()
}