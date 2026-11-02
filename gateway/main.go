package main

import (
	"log"
	"net/http"

	"example.com/finance/gateway/internal/api"
)

func main() {
	handler := api.NewRouter()

	addr := ":8080"
	log.Printf("Gateway service started on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}