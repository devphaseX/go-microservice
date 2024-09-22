package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (c *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler((cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders:   []string{"Accept", "Authorization", "X-CSRF-TOKEN"},
		AllowCredentials: true,
		MaxAge:           300,
	})))

	mux.Use(middleware.Heartbeat("ping"))

	mux.Post("/auth/sign-in", c.Signin)
	return mux
}
