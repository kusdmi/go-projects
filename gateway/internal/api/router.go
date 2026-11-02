package api

import "net/http"

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/transactions", CreateTransactionHandler)
	mux.HandleFunc("GET /api/transactions", ListTransactionsHandler)
	mux.HandleFunc("POST /api/budgets", CreateBudgetHandler)
	mux.HandleFunc("GET /api/budgets", ListBudgetsHandler)

	return LoggingMiddleware(mux)
}