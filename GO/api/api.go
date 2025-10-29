package api

// @title Subscriptions API
// @version 1.0
// @description REST API для управления онлайн-подписками пользователей
// @host localhost:8080
// @BasePath /

// @contact.name Artem
// @contact.email disaer21@yandex.ru

import (
	"github.com/Halturshik/EM-test-task/GO/database"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type API struct {
	Store *database.Store
}

func NewAPI(store *database.Store) *API {
	return &API{Store: store}
}

func (api *API) Init(r *chi.Mux) {
	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", api.CreateSubscriptionHandler)
	})

	r.Route("/users/{user_id}/subscriptions", func(r chi.Router) {
		r.Get("/", api.GetSubscriptionsHandler)
		r.Get("/{service_name}", api.GetSubscriptionsHandler)
		r.Put("/{service_name}", api.UpdateSubscriptionHandler)
		r.Delete("/{service_name}", api.DeleteSubscriptionHandler)
		r.Post("/{service_name}/total", api.GetTotalSubscriptionCostHandler)

	})

	r.Get("/swagger/*", httpSwagger.Handler())

}
