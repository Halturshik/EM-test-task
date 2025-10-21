package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *API) getSubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	if strings.TrimSpace(userIDStr) == "" {
		log.Println("Ошибка: не указан uuid пользователя")
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Не указан идентификатор пользователя"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Println("Ошибка: некорректный формат UUID:", err)
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Некорректный формат идентификатора"})
		return
	}

	serviceName := chi.URLParam(r, "service_name")

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

	resp := make([]SubsResponse, 0, len(subsFromDB))
	for _, s := range subsFromDB {
		var endStr *string
		if s.EndDate != nil {
			tmp := s.EndDate.Format("01-2006")
			endStr = &tmp
		}
		resp = append(resp, SubsResponse{
			ServiceName: s.ServiceName,
			Price:       s.Price,
			StartDate:   s.StartDate.Format("01-2006"),
			EndDate:     endStr,
		})
	}

	writeJSON(w, http.StatusOK, resp)
}
