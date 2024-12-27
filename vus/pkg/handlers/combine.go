package handlers

import (
	"common"
	"fmt"
	"net/http"
)

func (h *Handler) combineFile(w http.ResponseWriter, r *http.Request) error {
	var cbn common.CompleteMultiPartUpload
	err := common.ParseJSON(r, &cbn)
	if err != nil {
		return fmt.Errorf("error parsing request: %w", err)
	}

	err = h.storage.CombineMultiPartUploads(r.Context(), cbn)
	if err != nil {
		return err
	}

	return common.WriteJSON(w, http.StatusOK, map[string]string{"msg": "upload completed"})

}
