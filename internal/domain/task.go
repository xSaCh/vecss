package domain

type MqTask struct {
	UploadId    string `json:"upload_id"`
	Key         string `json:"key"`
	Url         string `json:"url"`
	Resolutions []int  `json:"resolutions"`
	Thumbnail   bool   `json:"thumbnail"`
}
