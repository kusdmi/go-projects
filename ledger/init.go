package ledger

import (
	"example.com/finance/ledger/internal/cache"
	"example.com/finance/ledger/internal/db"
)

// Init подключает базу данных и Redis. Вызывается из Gateway или cmd.
func Init() error {
	if err := db.Connect(); err != nil {
		return err
	}
	return cache.Connect()
}