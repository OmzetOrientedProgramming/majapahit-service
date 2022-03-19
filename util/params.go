package util

import "github.com/xendit/xendit-go"

var (
	// XenditFeesDefault Default params for additional invoices in xendit
	XenditFeesDefault = []xendit.InvoiceFee{
		{
			Type:  "Booking Fee",
			Value: 15000,
		},
		{
			Type:  "PlatformFee",
			Value: 3000,
		},
	}
)

const (
	// InvoiceDuration expired duration when creating invoice
	InvoiceDuration = 7200 // 2 hours

	// OOPEmail for Omzet Oriented Programming
	OOPEmail = "pplb.oop@gmail.com"

	// BankBCA for xendit
	BankBCA = "BCA"

	// DefaultLimit for Pagination
	DefaultLimit = 10

	// DefaultPage for Pagination
	DefaultPage = 1

	// MaxLimit for Pagination
	MaxLimit = 100
)

var (
	// DefaultPaymentMethod for xendit
	DefaultPaymentMethod = []string{"BCA", "OVO", "DANA", "QRIS"}

	// SMSNotification for xendit default notification
	SMSNotification = []string{"sms"}
)
