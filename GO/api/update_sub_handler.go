package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Halturshik/EM-test-task/GO/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// @Summary Обновить подписку
// @Description Обновляет уровень и/или дату окончания подписки пользователя
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id path string true "UUID пользователя"
// @Param service_name path string true "Название сервиса"
// @Param body body api.UpdateSubRequest true "Новые данные подписки"
// @Success 200 {object} api.UpdateSubResponse "Сообщение об обновлении подписки"
// @Failure 400 {object} api.ErrorResponse "Некорректные данные запроса"
// @Failure 404 {object} api.ErrorResponse "Подписка не найдена"
// @Failure 500 {object} api.ErrorResponse "Внутренняя ошибка сервера"
// @Router /users/{user_id}/subscriptions/{service_name} [put]
func (api *API) UpdateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	serviceName := strings.TrimSpace(chi.URLParam(r, "service_name"))

	if strings.TrimSpace(userIDStr) == "" || serviceName == "" {
		log.Println("Ошибка: не указан uuid пользователя или название сервиса")
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Не указан идентификатор пользователя или название сервиса подписки"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Println("Ошибка: некорректный формат uuid:", err)
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Некорректный формат идентификатора пользователя"})
		return
	}

	reSN := regexp.MustCompile(`^[A-Za-z0-9 ]+$`)
	if !reSN.MatchString(serviceName) {
		log.Println("Ошибка: в названии сервиса используются недопустимые символы")
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Недопустимое название сервиса: используйте только буквы, цифры и пробелы"})
		return
	}

	var req struct {
		NewPrice   *int    `json:"new_price,omitempty"`
		NewEndDate *string `json:"new_end_date,omitempty"`
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

	var newEndDateParsed *time.Time
	newEndDateProvided := false
	if req.NewEndDate != nil {
		if strings.TrimSpace(*req.NewEndDate) != "" {
			t, err := time.Parse("01-2006", *req.NewEndDate)
			if err != nil {
				log.Println("Ошибка: некорректный формат даты конца подписки")
				writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Неверный формат даты окончания действия подписки (используйте месяц-год)"})
				return
			}

			now := time.Now()
			currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			if t.Before(currentMonth) {
				log.Println("Ошибка: дата конца подписки в прошлом")
				writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Дата окончания подписки не может быть раньше текущего месяца"})
				return
			}
			endOfMonth := time.Date(t.Year(), t.Month()+1, 0, 23, 59, 59, 0, t.Location())
			newEndDateParsed = &endOfMonth
			newEndDateProvided = true
		}
	}

	priceChanged, endDateChanged, opType, err := api.Store.UpdateSubscription(r.Context(), userID, serviceName, req.NewPrice, newEndDateParsed, newEndDateProvided)
	if err != nil {
		if errors.Is(err, database.ErrSubNotFound) {
			log.Println("Ошибка: активная подписка не найдена")
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "Активная подписка не найдена"})
			return
		}

		log.Println("Ошибка при обновлении подписки:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "Не удалось обновить подписку. Повторите попытку позже"})
		return
	}

	if opType == "" && !priceChanged && !endDateChanged {
		log.Println("Ошибка: подписка уже соответствует поступившим параметрам")
		writeJSON(w, http.StatusOK, map[string]any{"message": "Выбранная подписка уже соответствует указанным параметрам"})
		return
	}

	var parts []string

	switch opType {
	case "upgrade":
		parts = append(parts, "Уровень подписки повышен и уже действует")
	case "downgrade":
		parts = append(parts, "Уровень подписки понижен, но вступит в силу в следующем месяце. До конца месяца сохраняется текущий уровень подписки")
	case "rollback":
		parts = append(parts, "Вернули прежний уровень подписки")
	}

	if endDateChanged {
		parts = append(parts, "Дата окончания подписки изменена")
	}

	msg := strings.Join(parts, ". ")

	writeJSON(w, http.StatusOK, map[string]any{"message": msg})
}
