package main

import (
	"log"
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("method not allowed"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("pong"))
}

func main() {
	http.HandleFunc("/ping", pingHandler)

	addr := ":8080"
	log.Printf("Gateway service started on %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
