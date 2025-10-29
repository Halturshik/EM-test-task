package api

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// @Summary Получить подписки пользователя
// @Description Возвращает список подписок для указанного user_id
// @Tags subscriptions
// @Produce json
// @Param user_id path string true "UUID пользователя"
// @Success 200 {array} api.SubResponse "Список подписок"
// @Failure 400 {object} api.ErrorResponse "Некорректный UUID пользователя или service_name"
// @Failure 500 {object} api.ErrorResponse "Внутренняя ошибка сервера"
// @Router /users/{user_id}/subscriptions [get]
// @Router /users/{user_id}/subscriptions/{service_name} [get]
func (api *API) GetSubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	serviceName := strings.TrimSpace(chi.URLParam(r, "service_name"))

	if strings.TrimSpace(userIDStr) == "" {
		log.Println("Ошибка: не указан uuid пользователя")
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Не указан идентификатор пользователя"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Println("Ошибка: некорректный формат uuid:", err)
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Некорректный формат идентификатора"})
		return
	}

	if serviceName != "" {
		reSN := regexp.MustCompile(`^[A-Za-z0-9 ]+$`)
		if !reSN.MatchString(serviceName) {
			log.Println("Ошибка: в названии сервиса используются недопустимые символы")
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Недопустимое название сервиса: используйте только буквы, цифры и пробелы"})
			return
		}
	}

	subsFromDB, err := api.Store.GetSubscriptions(r.Context(), userID, serviceName)
	if err != nil {
		log.Println("Ошибка: не удалось вытащить подписку", err)
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "Не удалось произвести поиск подписки. Повторите попытку позже"})
		return
	}

	if len(subsFromDB) == 0 {
		writeJSON(w, http.StatusOK, map[string]any{"message": "Подписок не найдено"})
		return
	}

	type subsResponse struct {
		ServiceName string  `json:"service_name"`
		Price       int     `json:"price"`
		StartDate   string  `json:"start_date"`
		EndDate     *string `json:"end_date,omitempty"`
	}

	resp := make([]subsResponse, 0, len(subsFromDB))
	for _, s := range subsFromDB {
		var endStr *string
		infiniteDate := time.Date(2099, 12, 31, 0, 0, 0, 0, s.EndDate.Location())
		if !s.EndDate.Equal(infiniteDate) {
			tmp := s.EndDate.Format("01-2006")
			endStr = &tmp
		}
		resp = append(resp, subsResponse{
			ServiceName: s.ServiceName,
			Price:       s.Price,
			StartDate:   s.StartDate.Format("01-2006"),
			EndDate:     endStr,
		})
	}

	writeJSON(w, http.StatusOK, resp)
}
