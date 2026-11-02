package ledger

import (
	"encoding/json"
	"fmt"
	"io"
)

type Budget struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
	Period   string  `json:"period,omitempty"`
}

func (b Budget) Validate() error {
	if b.Category == "" {
		return fmt.Errorf("budget category cannot be empty")
	}
	if b.Limit <= 0 {
		return fmt.Errorf("budget limit must be positive, got %.2f", b.Limit)
	}
	return nil
}

var budgets = map[string]Budget{}

func SetBudget(b Budget) error {
	if err := b.Validate(); err != nil {
		return err
	}
	budgets[b.Category] = b
	return nil
}

func GetBudget(category string) (Budget, bool) {
	b, ok := budgets[category]
	return b, ok
}

func LoadBudgets(r io.Reader) error {
	var items []Budget
	dec := json.NewDecoder(r)
	if err := dec.Decode(&items); err != nil {
		return fmt.Errorf("parse budgets JSON: %w", err)
	}

	for _, b := range items {
		if err := SetBudget(b); err != nil {
			return fmt.Errorf("set budget %q: %w", b.Category, err)
		}
	}

	return nil
}

func ListBudgets() []Budget {
	result := make([]Budget, 0, len(budgets))
	for _, b := range budgets {
		result = append(result, b)
	}
	return result
}