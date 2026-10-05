package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type Budget struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
	Period   string  `json:"period,omitempty"`
}

var budgets = map[string]Budget{}

func SetBudget(b Budget) error {
	if b.Category == "" {
		return errors.New("category cannot be empty")
	}
	if b.Limit <= 0 {
		return errors.New("limit must be positive")
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
