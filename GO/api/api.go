package api

import (
	"github.com/Halturshik/EM-test-task/GO/db"
	"github.com/go-chi/chi/v5"
)

type API struct {
	Store *db.Store
}

func NewAPI(store *db.Store) *API {
	return &API{Store: store}
}

func (api *API) Init(r *chi.Mux) {

}
