package item

type ListItem struct {
	Items []Item `json:"items"`
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