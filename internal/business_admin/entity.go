package businessadmin

// BalanceDetail consist related information for balance
type BalanceDetail struct {
	LatestDisbursementDate string  `json:"latest_disbursement_date"`
	Balance                float64 `json:"balance"`
}

// DisbursementDetail consist related information for disbursement
type DisbursementDetail struct {
	Date   string
	Amount float64
	Status int
}

// ListTransaction is a container for transaction history of customers
type ListTransaction struct {
	Transactions []Transaction `json:"transaction"`
	TotalCount   int           `json:"total_count"`
}

// Transaction consist related information for transaction history from customer
type Transaction struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Price int    `json:"price" db:"total_price"`
	Date  string `json:"date"`
}

// ListTransactionRequest consists of request data from client
type ListTransactionRequest struct {
	Limit  int    `json:"limit"`
	Page   int    `json:"page"`
	Path   string `json:"path"`
	UserID int    `json:"user_id"`
}

// TransactionHistoryDetail consists detail of a transaction history from customer
type TransactionHistoryDetail struct {
	CustomerName   string       `json:"customer_name"`
	CustomerImage  string       `json:"customer_image"`
	Date           string       `json:"date"`
	StartTime      string       `db:"start_time" json:"start_time"`
	EndTime        string       `db:"end_time" json:"end_time"`
	Capacity       int          `json:"capacity"`
	TotalPriceItem float64      `json:"total_price_item" db:"total_price"`
	Items          []ItemDetail `json:"items"`
}

// CustomerForTrasactionHistoryDetail is a wrapper consists customer information
type CustomerForTrasactionHistoryDetail struct {
	CustomerName  string `json:"customer_name" db:"name"`
	CustomerImage string `json:"customer_image" db:"image"`
}

// ItemsWrapper is a struct to provide a wrap for items
type ItemsWrapper struct {
	Items []ItemDetail
}

// ItemDetail consist information related to item
type ItemDetail struct {
	Name  string  `json:"name"`
	Qty   int     `json:"qty"`
	Price float64 `json:"price"`
}
