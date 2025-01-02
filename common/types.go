package common

import "time"

type MultiPartUrls struct {
	UploadId  string    `json:"upload_id"`
	Urls      []string  `json:"urls"`
	ChunkSize int       `json:"chunk_size"`
	CreateAt  time.Time `json:"create_at"`
	ExpireAt  time.Time `json:"expire_at"`
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

type MqTask struct {
	UploadId    string   `json:"upload_id"`
	Key         string   `json:"key"`
	Resolutions []string `json:"resolutions"`
	Thumbnail   bool     `json:"thumbnail"`
}
