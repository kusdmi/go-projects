package ledger

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"example.com/finance/ledger/internal/cache"
	"example.com/finance/ledger/internal/db"
)

// GetReportSummary возвращает суммы расходов по категориям за период (с кэшированием в Redis)
func GetReportSummary(ctx context.Context, from, to string) (map[string]float64, error) {
	cacheKey := fmt.Sprintf("report:summary:%s:%s", from, to)

	// 1. Пытаемся достать из Redis
	if val, err := cache.Client.Get(ctx, cacheKey).Result(); err == nil {
		var summary map[string]float64
		if err := json.Unmarshal([]byte(val), &summary); err == nil {
			fmt.Println("Cache HIT for key:", cacheKey)
			return summary, nil
		}
	}

	fmt.Println("Cache MISS for key:", cacheKey)

	// 2. Идем в БД
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

	// 3. Сохраняем в Redis на 30 секунд
	if data, err := json.Marshal(summary); err == nil {
		cache.Client.Set(ctx, cacheKey, data, 30*time.Second)
	}

	return summary, nil
}