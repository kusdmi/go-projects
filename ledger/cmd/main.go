package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"example.com/finance/ledger"
	"example.com/finance/ledger/internal/cache"
	"example.com/finance/ledger/internal/db"
)

func main() {
	fmt.Println("Ledger service started")

	if err := db.Connect(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	if err := cache.Connect(); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	file, err := os.Open("budgets.json")
	if err != nil {
		log.Fatalf("failed to open budgets.json: %v", err)
	}
	defer file.Close()

	// Демонстрация загрузки бюджетов (опционально, можно удалить)
	// ... (код для LoadBudgets, если он нужен, но он больше не использует БД)

	ctx := context.Background()
	from := "2025-09-01"
	to := "2025-09-30"
	summary, err := GetReportSummary(ctx, from, to)
	if err != nil {
		log.Printf("Report summary error: %v", err)
	} else {
		fmt.Printf("Report summary: %+v\n", summary)
	}

	// Демонстрация: повторный вызов должен быть быстрее (из кэша)
	start := time.Now()
	_, _ = GetReportSummary(ctx, from, to)
	fmt.Printf("Second call took %v (cached)\n", time.Since(start))
}

// GetReportSummary — пример агрегации с кэшированием
func GetReportSummary(ctx context.Context, from, to string) (map[string]float64, error) {
	cacheKey := fmt.Sprintf("report:summary:%s:%s", from, to)

	if val, err := cache.Client.Get(ctx, cacheKey).Result(); err == nil {
		var summary map[string]float64
		if err := json.Unmarshal([]byte(val), &summary); err == nil {
			return summary, nil
		}
	}

	query := `
		SELECT category, SUM(amount)
		FROM expenses
		WHERE date BETWEEN $1 AND $2
		GROUP BY category
	`
	rows, err := db.DB.Query(query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summary := make(map[string]float64)
	for rows.Next() {
		var category string
		var total float64
		if err := rows.Scan(&category, &total); err != nil {
			return nil, err
		}
		summary[category] = total
	}

	data, _ := json.Marshal(summary)
	cache.Client.Set(ctx, cacheKey, data, 30*time.Second)

	return summary, nil
}