package api

import (
	"github.com/Halturshik/EM-test-task/GO/database"
	"github.com/go-chi/chi/v5"
)

type SubsResponse struct {
	ID          int64   `json:"id"`
	UserID      string  `json:"user_id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

type API struct {
	Store *database.Store
}

func NewAPI(store *database.Store) *API {
	return &API{Store: store}
}

func (api *API) Init(r *chi.Mux) {
	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", api.createSubscriptionHandler)
	})

	r.Route("/users/{user_id}/subscriptions", func(r chi.Router) {
		r.Get("/", api.getSubscriptionsHandler)
		r.Get("/{service_name}", api.getSubscriptionsHandler)
		r.Put("/{service_name}", api.updateSubscriptionHandler)
	})

}
