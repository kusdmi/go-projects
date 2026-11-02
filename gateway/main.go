package main

import (
	"log"
	"net/http"

	"example.com/finance/gateway/internal/api"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/transactions", api.CreateTransactionHandler)
	mux.HandleFunc("GET /api/transactions", api.ListTransactionsHandler)
	mux.HandleFunc("POST /api/budgets", api.CreateBudgetHandler)
	mux.HandleFunc("GET /api/budgets", api.ListBudgetsHandler)

	handler := api.LoggingMiddleware(mux)

	addr := ":8080"
	log.Printf("Gateway service started on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
