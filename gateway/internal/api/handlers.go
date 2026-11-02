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

	// Передаем указатель &tx, чтобы Ledger сохранил в него ID из базы данных (RETURNING id)
	if err := ledger.AddTransaction(&tx); err != nil {
		if strings.Contains(err.Error(), "budget exceeded") {
			writeError(w, http.StatusConflict, "budget exceeded")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Теперь tx содержит сгенерированный ID, возвращаем его клиенту
	writeJSON(w, http.StatusCreated, FromDomainTransaction(tx))
}

func ListTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	// Теперь ListTransactions возвращает ошибку
	txs, err := ledger.ListTransactions()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list transactions")
		return
	}

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
	// Теперь ListBudgets возвращает ошибку
	bs, err := ledger.ListBudgets()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list budgets")
		return
	}

	res := make([]BudgetResponse, 0, len(bs))
	for _, b := range bs {
		res = append(res, FromDomainBudget(b))
	}
	writeJSON(w, http.StatusOK, res)
}

func GetReportSummaryHandler(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	if from == "" || to == "" {
		writeError(w, http.StatusBadRequest, "from and to query parameters are required")
		return
	}

	summary, err := ledger.GetReportSummary(r.Context(), from, to)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get report summary")
		return
	}

	writeJSON(w, http.StatusOK, summary)
}