package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Halturshik/EM-test-task/GO/database"
	"github.com/google/uuid"
)

func (api *API) createSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID      string  `json:"user_id"`
		ServiceName string  `json:"service_name"`
		Price       int     `json:"price"`
		StartDate   string  `json:"start_date"`
		EndDate     *string `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Ошибка: не удалось прочитать тело запроса", err)
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Некорректно оформлено тело запроса"})
		return
	}

	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Println("Ошибка: некорректный формат uuid")
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Указан некорректный формат идентификатора пользователя"})
		return
	}

	if strings.TrimSpace(req.ServiceName) == "" {
		log.Println("Ошибка: не указан сервис подписки")
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Не указан сервис подписки"})
		return
	}

	validPrices := map[int]bool{50: true, 100: true, 200: true}
	if !validPrices[req.Price] {
		log.Println("Ошибка: выбран несуществующий уровень подписки")
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Выберите допустимый уровень подписки: Базовый (50), Продвинутый (100), Премиум (200)"})
		return
	}

	start, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		log.Println("Ошибка: некорректный формат даты начала")
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Неверный формат даты начала действия подписки (используйте месяц-год)"})
		return
	}

	var end *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		endParsed, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			log.Println("Ошибка: некорректный формат даты начала")
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Неверный формат даты окончания действия подписки (используйте месяц-год)"})
			return
		}

		if endParsed.Before(start) {
			log.Println("Ошибка: дата окончания подписки раньше даты начала")
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Дата окончания действия подписки не может быть раньше даты ее начала действия"})
			return
		}
		end = &endParsed
	}
	sub := &database.Subs{
		UserID:      uid,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   start,
		EndDate:     end,
	}

	err = api.Store.CreateSubscription(r.Context(), sub)
	if err != nil {
		if errors.Is(err, database.ErrSubIsExist) {
			log.Println("Ошибка: подписка уже существует")
			writeJSON(w, http.StatusConflict, map[string]any{"error": "Подписка на выбранный сервис уже существует"})
			return
		}

		log.Println("Ошибка: не удалось создать подписку", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "Не удалось создать подписку. Повторите попытку позже"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"message": "Подписка успешно создана"})

}
