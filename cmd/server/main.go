package main

import (
	"net/http"

	configs "github.com/gabriel-ulisses-andrade/goexpert-desafio-multithreading-api/configs"
	"github.com/gabriel-ulisses-andrade/goexpert-desafio-multithreading-api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	cepHandler := handlers.NewCEPHandler(
		configs.URLBrasilApi,
		configs.URLViaCep,
	)

	r := chi.NewRouter()

	r.Get("/cep/{cep}", cepHandler.ConsultarCep)

	http.ListenAndServe(":8080", r)
}
