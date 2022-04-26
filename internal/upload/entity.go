package upload

type FileRequest struct {
	File         string `json:"file"`
	CustomerName string
}

type FileResponse struct {
	URL string `json:"url"`
}
