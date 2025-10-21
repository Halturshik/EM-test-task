package main

import (
	"log"
	"net/http"

	"github.com/Halturshik/EM-test-task/GO/api"
	"github.com/Halturshik/EM-test-task/GO/database"
	"github.com/Halturshik/EM-test-task/config"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Предупреждение: .env файл не найден, будут использоваться переменные окружения")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	dbConnection, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Ошибка при подключении к БД:", err)
	}
	defer dbConnection.Close()

	store := database.NewStore(dbConnection)
	apiServer := api.NewAPI(store)

	r := chi.NewRouter()

	r.Use(api.LoggingMiddleware)

	apiServer.Init(r)

	log.Printf("Сервер запущен и слушает порт %s\n", cfg.AppPort)
	if err := http.ListenAndServe(":"+cfg.AppPort, r); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}
