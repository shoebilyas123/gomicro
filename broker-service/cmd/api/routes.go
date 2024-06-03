package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// Setup the cors config
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://*"},
		AllowedMethods: []string{"GET","POST","PUT","PATCH","DELETE"},
		AllowedHeaders: []string{"Accept","Authorization","Conten-Type","X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge: 300,
	}))

	mux.Use(middleware.Heartbeat("/ping"));

	mux.Post("/", app.Broker)
	return mux;
}