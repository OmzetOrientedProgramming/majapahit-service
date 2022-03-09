package checkup

// APIResponse for general API response for this package
type APIResponse struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    *interface{} `json:"data,omitempty"`
	Errors  []string     `json:"errors,omitempty"`
}
