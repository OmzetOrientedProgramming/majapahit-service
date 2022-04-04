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
