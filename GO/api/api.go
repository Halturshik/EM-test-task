package api

type API struct {
	Store *db.Store
}

func NewAPI(store *db.Store) *API {
	return &API{Store: store}
}

func (api *API) Init(r *chi.Mux) {

}
