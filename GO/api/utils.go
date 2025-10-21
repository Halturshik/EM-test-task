package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := json.Marshal(data)
	if err != nil {
		log.Println("Ошибка при формировании JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		resp = []byte(`{"error":"ошибка при формировании JSON"}`)
	} else {
		w.WriteHeader(status)
	}

	if _, err := w.Write(resp); err != nil {
		log.Println("Ошибка при отправке ответа:", err)
	}
}
