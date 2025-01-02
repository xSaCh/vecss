package handlers

import (
	"common"
	"fmt"
	"net/http"

	"common/mq"

	"github.com/gorilla/mux"
	"github.com/xSaCh/vecss/vus/pkg/repositories"
)

const CHUNK_SIZE = 15 * 1024 * 1024

type Handler struct {
	storage repositories.Storage
	emitter mq.Emitter
}

func NewHandler(storage repositories.Storage, emitter mq.Emitter) *Handler {
	return &Handler{storage: storage, emitter: emitter}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/upload/", common.MakeHTTPHandleFunc(h.uploadGet)).Methods(http.MethodGet)
	router.HandleFunc("/upload", common.MakeHTTPHandleFunc(h.uploadGet)).Methods(http.MethodGet)

	router.HandleFunc("/upload/", common.MakeHTTPHandleFunc(h.uploadFile)).Methods(http.MethodPost)
	router.HandleFunc("/upload", common.MakeHTTPHandleFunc(h.uploadFile)).Methods(http.MethodPost)

	router.HandleFunc("/combine/", common.MakeHTTPHandleFunc(h.combineFile)).Methods(http.MethodPost)
	router.HandleFunc("/combine", common.MakeHTTPHandleFunc(h.combineFile)).Methods(http.MethodPost)
}

func (h *Handler) uploadFile(w http.ResponseWriter, r *http.Request) error {
	file, hdr, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()
	if hdr.Size == 0 {
		return fmt.Errorf("file size is 0")
	}

	numofParts := (hdr.Size / CHUNK_SIZE) + 1
	parts := make([]int, numofParts)
	for i := 0; i < int(numofParts); i++ {
		parts[i] = i + 1
	}

	// TODO: Generate unqiue key for the file
	key := hdr.Filename

	urls, err := h.storage.GenerateMultiPartPreSignedUrls(r.Context(), key, parts)
	if err != nil {
		return fmt.Errorf("error generating pre-signed urls: %w", err)
	}
	urls.ChunkSize = CHUNK_SIZE
	return common.WriteJSON(w, http.StatusOK, urls)

}

func (h *Handler) uploadGet(w http.ResponseWriter, r *http.Request) error {
	fmt.Printf("r.Header: %v\n", r.Header)
	return common.WriteJSON(w, http.StatusOK, map[string]string{"msg": "use post method to upload file"})

}
