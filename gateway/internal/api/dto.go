package api

import "example.com/finance/ledger"

type CreateTransactionRequest struct {
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

type TransactionResponse struct {
	ID          int     `json:"id"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

type CreateBudgetRequest struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
}

type BudgetResponse struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
	Period   string  `json:"period,omitempty"`
}

func (r CreateTransactionRequest) ToDomain() ledger.Transaction {
	return ledger.Transaction{
		Amount:      r.Amount,
		Category:    r.Category,
		Description: r.Description,
		Date:        r.Date,
	}
}

func FromDomainTransaction(tx ledger.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:          tx.ID,
		Amount:      tx.Amount,
		Category:    tx.Category,
		Description: tx.Description,
		Date:        tx.Date,
	}
}

func (r CreateBudgetRequest) ToDomain() ledger.Budget {
	return ledger.Budget{
		Category: r.Category,
		Limit:    r.Limit,
	}
}

func FromDomainBudget(b ledger.Budget) BudgetResponse {
	return BudgetResponse{
		Category: b.Category,
		Limit:    b.Limit,
		Period:   b.Period,
	}
}
