package pkg

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xSaCh/vecss/vus/pkg/handlers"
	"github.com/xSaCh/vecss/vus/pkg/mq"
	"github.com/xSaCh/vecss/vus/pkg/repositories"
)

type APIServer struct {
	addr    string
	storage repositories.Storage
	emitter mq.Emitter
}

func NewAPIServer(addr string, storage repositories.Storage, emitter mq.Emitter) *APIServer {
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
