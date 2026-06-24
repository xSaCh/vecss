package pkg

import (
	"log"
	"net/http"

	"vecss/internal/mq"
	"vecss/internal/storage"

	"vecss/internal/vus/pkg/handlers"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr    string
	storage storage.Storage
	emitter mq.Emitter
}

func NewAPIServer(addr string, storage storage.Storage, emitter mq.Emitter) *APIServer {
	return &APIServer{
		addr:    addr,
		storage: storage,
		emitter: emitter,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	uploadSubRouter := router.PathPrefix("/").Subrouter()
	handlers.NewHandler(s.storage, s.emitter).RegisterRoutes(uploadSubRouter)

	router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(`{"message": "okie"}`))
	})

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
