package main

import (
	"log"
	"net/http"

	"github.com/Halturshik/EM-test-task/GO/api"
	"github.com/Halturshik/EM-test-task/GO/db"
	"github.com/Halturshik/EM-test-task/config"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	dbConnection, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Ошибка при подключении к БД:", err)
	}
	defer dbConnection.Close()

	store := db.NewStore(dbConnection)
	apiServer := api.NewAPI(store)

	r := chi.NewRouter()

	apiServer.Init(r)

	log.Println("Сервер запущен на порту", cfg.AppPort)
	if err := http.ListenAndServe(":"+cfg.AppPort, r); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}
