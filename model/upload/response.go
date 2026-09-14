package upload

type UploadImageResponse struct {
	URL       string `json:"url"`
	ObjectKey string `json:"object_key"`
	Scope     string `json:"scope"`
	Resource  string `json:"resource_id"`
}