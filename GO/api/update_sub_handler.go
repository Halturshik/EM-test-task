package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Halturshik/EM-test-task/GO/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *API) updateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	serviceName := chi.URLParam(r, "service_name")

	if strings.TrimSpace(userIDStr) == "" || strings.TrimSpace(serviceName) == "" {
		log.Println("Ошибка: не указан uuid пользователя или название сервиса")
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Не указан идентификатор пользователя или название сервиса подписки"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Println("Ошибка: некорректный формат UUID:", err)
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Некорректный формат идентификатора пользователя"})
		return
	}

	var req struct {
		NewPrice   *int       `json:"new_price,omitempty"`
		NewEndDate *time.Time `json:"new_end_date,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Ошибка: не удалось прочитать тело запроса", err)
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Некорректно оформлено тело запроса"})
		return
	}

	if req.NewPrice == nil && req.NewEndDate == nil {
		log.Println("Ошибка: не указаны поля для изменения")
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Не заполнены поля для обновления"})
		return
	}

	validPrices := map[int]bool{50: true, 100: true, 200: true}
	if req.NewPrice != nil && !validPrices[*req.NewPrice] {
		log.Println("Ошибка: выбран несуществующий уровень подписки")
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Выберите допустимый уровень подписки: Базовый (50), Продвинутый (100), Премиум (200)"})
		return
	}

	err = api.Store.UpdateSubscription(r.Context(), userID, serviceName, req.NewPrice, req.NewEndDate)
	if err != nil {
		if errors.Is(err, database.ErrSubNotFound) {
			log.Println("Ошибка: подписка не найдена")
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "Подписка не найдена"})
			return
		}

		log.Println("Ошибка при обновлении подписки:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "Не удалось обновить подписку. Повторите попытку позже"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"message": "Подписка успешно обновлена"})
}
