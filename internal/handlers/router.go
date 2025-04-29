package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/"))

	r.Route("/v1", func(v1 chi.Router) {
		v1.Route("/swift-codes", func(swift chi.Router) {
			swift.Get("/{swiftCode}", GetBankDataBySwift)
			swift.Get("/country/{countryISO2code}", GetBanksDataByCountryISO2)
			swift.Post("/", CreateBank)
			swift.Delete("/{swiftCode}", DeleteSwiftCode)
		})
	})

	return r
}
