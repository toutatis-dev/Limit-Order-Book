package main

import (
	"net/http"

	"lob/internal/api"
	"lob/internal/engine"

	"github.com/go-chi/chi/v5"
)

func main() {
	e := engine.NewEngine()
	a := api.NewAPI(e)

	r := chi.NewRouter()

	r.Post("/order", a.HandleCreateOrder)
	r.Delete("/order/{id}", a.HandleCancelOrder)
	r.Get("/book", a.HandleGetBook)

	http.ListenAndServe(":8080", r)
}
