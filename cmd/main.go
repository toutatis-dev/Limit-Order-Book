package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"lob/internal/api"
	"lob/internal/engine"
	"lob/internal/store"

	"github.com/go-chi/chi/v5"
)

type EnvVars struct {
	DBUser string
	DBPort string
	DBPass string
	DBHost string
	DBName string
}

func main() {
	env := EnvVars{
		DBUser: os.Getenv("DB_USER"),
		DBPort: os.Getenv("DB_PORT"),
		DBPass: os.Getenv("DB_PASS"),
		DBHost: os.Getenv("DB_HOST"),
		DBName: os.Getenv("DB_NAME"),
	}
	pgUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", env.DBUser, env.DBPass, env.DBHost, env.DBPort, env.DBName)
	store, err := store.NewStore(pgUrl)
	if err != nil {
		log.Fatalf("Error creating NewStore: %v", err)
	}
	defer store.Close()

	e := engine.NewEngine(store)
	a := api.NewAPI(e)

	r := chi.NewRouter()

	r.Post("/order", a.HandleCreateOrder)
	r.Delete("/order/{id}", a.HandleCancelOrder)
	r.Get("/book", a.HandleGetBook)

	log.Fatal(http.ListenAndServe(":1770", r))
}
