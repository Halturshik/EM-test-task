package main

import (
	"log"
	"os"
)

func main() {

	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		dbPath = "subscribers.db"
	}

	dbConnection, err := db.ConnectDB(dbPath)
	if err != nil {
		log.Fatal("Ошибка при подключении к БД: %v", err)
	}
	defer dbConnection.Close()

	store := db.NewStore(dbConnection)
	apiServer := api.NewAPI(store)

	r := chi.NewRouter

	apiServer.Init(r)

}
