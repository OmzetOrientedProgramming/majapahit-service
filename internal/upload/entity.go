package upload

// FileRequest for upload file request struct
type FileRequest struct {
	File         string `json:"file"`
	CustomerName string
}

// FileResponse for upload file response struct
type FileResponse struct {
	URL string `json:"url"`
}
