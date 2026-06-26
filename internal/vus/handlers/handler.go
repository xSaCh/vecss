package handlers

import (
	"vecss/internal/common"
	"vecss/internal/mq"
	"vecss/internal/storage"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	storage storage.Storage
	emitter mq.Emitter
}

func NewHandler(storage storage.Storage, emitter mq.Emitter) *Handler {
	return &Handler{storage: storage, emitter: emitter}
}

func (h *Handler) RegisterUploadRoutes(r chi.Router) {
	r.Get("/upload", common.MakeHTTPHandleFunc(h.uploadGet))
	r.Post("/upload", common.MakeHTTPHandleFunc(h.uploadFile))
}

func (h *Handler) RegisterCombineRoute(r chi.Router) {
	r.Post("/combine", common.MakeHTTPHandleFunc(h.combineFile))
}
