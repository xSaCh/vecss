package common

import "time"

type MultiPartUrls struct {
	UploadId string    `json:"upload_id"`
	Urls     []string  `json:"urls"`
	CreateAt time.Time `json:"create_at"`
	ExpireAt time.Time `json:"expire_at"`
}

type CompleteMultiPartUpload struct {
	UploadId    string   `json:"upload_id"`
	ETags       []string `json:"etags"`
	PartNumbers []int    `json:"part_numbers"`
	Key         string   `json:"key"`
}

type ReqUploadFile struct {
	FileName string `json:"file_name"`
}
