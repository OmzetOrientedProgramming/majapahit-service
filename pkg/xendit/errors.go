package xendit

import "github.com/pkg/errors"

var (
	// ErrXenditCreateInvoice when calling create invoice
	ErrXenditCreateInvoice = errors.New("Xendit error create invoice")

	// ErrXenditCreateDisbursement when calling create disbursement
	ErrXenditCreateDisbursement = errors.New("Xendit error create disbursement")

	// ErrXenditGetInvoice when calling get invoice
	ErrXenditGetInvoice = errors.New("Xendit error get invoice")

	// ErrXenditGetDisbursement when calling get disbursement
	ErrXenditGetDisbursement = errors.New("Xendit error get disbursement")
)
