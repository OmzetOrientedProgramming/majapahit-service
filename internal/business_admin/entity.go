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
