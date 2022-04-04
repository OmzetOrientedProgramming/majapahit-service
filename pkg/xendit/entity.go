package xendit

// CreateDisbursementParams for disbursement params
type CreateDisbursementParams struct {
	ID                int      `json:"id"`
	BankAccountName   string   `json:"bank_acount_name"`
	BankAccountNumber string   `json:"bank_acount_number"`
	Amount            float64  `json:"amount"`
	Description       string   `json:"description"`
	Email             []string `json:"email"`
}

// CreateInvoiceParams for create invoices params
type CreateInvoiceParams struct {
	PlaceID             int     `json:"place_id"`
	Items               []Item  `json:"items"`
	Description         string  `json:"description"`
	CustomerName        string  `json:"customer_name"`
	CustomerPhoneNumber string  `json:"customer_phone_number"`
	BookingFee          float64 `json:"booking_fee"`
}

// Item that will be in invoice
type Item struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}
