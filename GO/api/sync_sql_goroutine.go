package api

import (
	"context"
	"log"
	"time"

	"github.com/Halturshik/EM-test-task/GO/database"
)

func StartMonthlySync(store *database.Store) {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
			duration := time.Until(next)

			time.Sleep(duration)

			if err := store.SyncSubscriptionPrices(context.Background()); err != nil {
				log.Println("Ошибка синхронизации подписок:", err)
			} else {
				log.Println("Синхронизация подписок выполнена успешно")
			}
		}
	}()
}
