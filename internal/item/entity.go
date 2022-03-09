package item

type ListItem struct {
	Items      []Item `json:"items"`
	TotalCount int    `json:"total_count"`
}

type Item struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type ListItemRequest struct {
	Limit   int    `json:"limit"`
	Page    int    `json:"page"`
	Path    string `json:"path"`
	PlaceID int    `json:"place_id"`
	Name    string `json:"name"`
}