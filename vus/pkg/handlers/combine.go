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

	url, err := h.storage.GetObjecPresigntUrl(r.Context(), cbn.Key)
	if err != nil {
		return err
	}
	//TODO: Add Meta Data to DB
	//TODO: Create Task for messageQueue
	task := common.MqTask{
		UploadId:    cbn.UploadId,
		Key:         cbn.Key,
		Url:         url,
		Resolutions: []int{1080, 720, 480},
		Thumbnail:   true,
	}
	err = h.emitter.Push(r.Context(), task)

	if err != nil {
		return err
	}
	return common.WriteJSON(w, http.StatusOK, map[string]string{"msg": "upload completed", "status": "pending"})

}
