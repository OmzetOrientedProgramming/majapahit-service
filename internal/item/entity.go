package item

// ListItem will be used as a container for items
type ListItem struct {
	Items      []Item `json:"items"`
	TotalCount int    `json:"total_count"`
}

// Item contains information that needed fot catalog items
type Item struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// ListItemRequest consists of request data from client
type ListItemRequest struct {
	Limit   int    `json:"limit"`
	Page    int    `json:"page"`
	Path    string `json:"path"`
	PlaceID int    `json:"place_id"`
	Name    string `json:"name"`
}
