package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"example.com/finance/ledger"
)

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	tx := req.ToDomain()

	if err := tx.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := ledger.AddTransaction(tx); err != nil {
		if strings.Contains(err.Error(), "budget exceeded") {
			writeError(w, http.StatusConflict, "budget exceeded")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	all := ledger.ListTransactions()
	created := all[len(all)-1]

	writeJSON(w, http.StatusCreated, FromDomainTransaction(created))
}

func ListTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	txs := ledger.ListTransactions()
	res := make([]TransactionResponse, 0, len(txs))
	for _, tx := range txs {
		res = append(res, FromDomainTransaction(tx))
	}
	writeJSON(w, http.StatusOK, res)
}

func CreateBudgetHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	b := req.ToDomain()
	if err := b.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := ledger.SetBudget(b); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, FromDomainBudget(b))
}

func ListBudgetsHandler(w http.ResponseWriter, r *http.Request) {
	bs := ledger.ListBudgets()
	res := make([]BudgetResponse, 0, len(bs))
	for _, b := range bs {
		res = append(res, FromDomainBudget(b))
	}
	writeJSON(w, http.StatusOK, res)
}
