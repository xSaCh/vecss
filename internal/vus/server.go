package vus

import (
	"net/http"

	"vecss/internal/mq"
	"vecss/internal/storage"
	"vecss/internal/vus/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRouters(storage storage.Storage, emitter mq.Emitter) http.Handler {
	router := chi.NewRouter()

	// Middlewares
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(res, req)
		})
	})

	// Routes
	handler := handlers.NewHandler(storage, emitter)
	router.Route("/", func(r chi.Router) {
		r.Get("/", func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(`{"message": "okie"}`))
		})
		handler.RegisterUploadRoutes(r)
		handler.RegisterCombineRoute(r)
	})
	return router
}
