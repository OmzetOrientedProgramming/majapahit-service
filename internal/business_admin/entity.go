package businessadmin

import "time"

// BalanceDetail consist related information for balance
type BalanceDetail struct {
	LatestDisbursementDate string  `json:"latest_disbursement_date"`
	Balance                float64 `json:"balance"`
}

// DisbursementDetail consist related information for disbursement
type DisbursementDetail struct {
	ID       int       `json:"id"`
	PlaceID  int       `json:"place_id"`
	Date     time.Time `json:"date"`
	XenditID string    `json:"xendit_id" db:"xendit_id"`
	Amount   float64   `json:"amount"`
	Status   int       `json:"status"`
}

// ListTransaction is a container for transaction history of customers
type ListTransaction struct {
	Transactions []Transaction `json:"transaction"`
	TotalCount   int           `json:"total_count"`
}

// Transaction consist related information for transaction history from customer
type Transaction struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Image string  `json:"image"`
	Price float64 `json:"price" db:"total_price"`
	Date  string  `json:"date"`
}

// ListTransactionRequest consists of request data from client
type ListTransactionRequest struct {
	Limit  int    `json:"limit"`
	Page   int    `json:"page"`
	Path   string `json:"path"`
	UserID int    `json:"user_id"`
}

// CreateDisbursementResponse for create disbursement response entity
type CreateDisbursementResponse struct {
	ID        int       `json:"place_id"`
	CreatedAt time.Time `json:"created_at"`
	Amount    float64   `json:"amount"`
	XenditID  string    `json:"xendit_id"`
}

// InfoForDisbursement for business admin info for disbursement entity
type InfoForDisbursement struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	BankAccountName   string `json:"bank_account_name" db:"bank_account_name"`
	BankAccountNumber string `json:"bank_account_number" db:"bank_account_number"`
	PlaceID           int    `json:"place_id" db:"place_id"`
}

// DisbursementCallback for disbursement callback struct
type DisbursementCallback struct {
	ID                      string  `json:"id"`
	ExternalID              string  `json:"external_id"`
	Amount                  float64 `json:"amount"`
	BankCode                string  `json:"bank_code"`
	AccountHolderName       string  `json:"account_holder_name"`
	DisbursementDescription string  `json:"disbursement_description"`
	FailureCode             string  `json:"failure_code"`
	Status                  string  `json:"status"`
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

// EditProfileRequest consist newest profile information about places
type EditProfileRequest struct {
	UserID      int
	Name        string `json:"name"`
	Description string `json:"description"`
}

// PlaceDetail contain important information in Place
type PlaceDetail struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name"`
	Image              string  `json:"image"`
	Address            string  `json:"address"`
	Description        string  `json:"description"`
	OpenHour           string  `json:"open_hour" db:"open_hour"`
	CloseHour          string  `json:"close_hour" db:"close_hour"`
	BookingPrice       int     `json:"booking_price" db:"booking_price"`
	MinSlot            int     `json:"min_slot" db:"min_slot_booking"`
	MaxSlot            int     `json:"max_slot" db:"max_slot_booking"`
	Capacity           int     `json:"capacity" db:"capacity"`
	MinIntervalBooking int     `json:"min_interval_booking" db:"min_interval_booking"`
	MaxIntervalBooking int     `json:"max_interval_booking" db:"max_interval_booking"`
	AverageRating      float64 `json:"average_rating" db:"rating"`
}

// ListReviewRequest consist of request for pagination and sorting purpose
type ListReviewRequest struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Path  string `json:"path"`
}
