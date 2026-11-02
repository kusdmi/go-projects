package main

import (
	"log"
	"net/http"

	"example.com/finance/gateway/internal/api"
	"example.com/finance/ledger"
)

func main() {
	// Инициализация БД и Redis
	if err := ledger.Init(); err != nil {
		log.Fatalf("Failed to initialize Ledger dependencies: %v", err)
	}

	handler := api.NewRouter()

	addr := ":8080"
	log.Printf("Gateway service started on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}