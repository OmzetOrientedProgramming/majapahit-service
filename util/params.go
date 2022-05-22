package util

import (
	"regexp"

	"github.com/xendit/xendit-go"
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

	// ApplicationJSON for content-type
	ApplicationJSON = "application/json"

	// GenderUndefined mapping for gender undefined
	GenderUndefined = 0

	// MinimumNameLength for user's minimum name length
	MinimumNameLength = 3

	// GenderMale mapping for gender male
	GenderMale = 1

	// GenderFemale mapping for gender female
	GenderFemale = 2

	// StatusCustomer for mapping status customer
	StatusCustomer = 0

	// StatusBusinessAdmin for mapping status business admin
	StatusBusinessAdmin = 1

	// TimeLayout for time layout convention
	TimeLayout = "15:04:05"
	// DateLayout for date layout convention
	DateLayout = "2006-01-02"

	// BookingMenungguKonfirmasi integer mapping
	BookingMenungguKonfirmasi = 0
	// BookingBelumMembayar integer mapping
	BookingBelumMembayar = 1
	// BookingBerhasil integer mapping
	BookingBerhasil = 2
	// BookingSelesai integer mapping
	BookingSelesai = 3
	// BookingGagal integer mapping
	BookingGagal = 4

	// Available booking status
	Available = 0
	//FullyBook booking status
	FullyBook = 1

	// XenditDisbursementPending for xendit disbursement status pending
	XenditDisbursementPending = 0
	// XenditDisbursementCompleted for xendit disbursement status completed
	XenditDisbursementCompleted = 1
	// XenditDisbursementFailed for xendit disbursemnet status failed
	XenditDisbursementFailed = 2

	// XenditDisbursementPendingString for xendit disbursement callback pending
	XenditDisbursementPendingString = "PENDING"
	// XenditDisbursementCompletedString for xendit disbursement callback completed
	XenditDisbursementCompletedString = "COMPLETED"
	// XenditDisbursementFailedString for xendit disbursement callback failed
	XenditDisbursementFailedString = "FAILED"

	// MaximumReviewLength for review content validation
	MaximumReviewLength = 500
	// MinimumRatingValue for rating valie validation
	MinimumRatingValue = 1
	// MaximumRatingValue for rating valie validation
	MaximumRatingValue = 5
)

var (
	// PhoneNumberRegex use for phone number validation
	PhoneNumberRegex = regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)

	// XenditFeesDefault Default params for additional invoices in xendit
	XenditFeesDefault = []xendit.InvoiceFee{
		{
			Type:  "PlatformFee",
			Value: XenditPlatformFee,
		},
	}

	// DefaultPaymentMethod for xendit
	DefaultPaymentMethod = []string{"BCA", "OVO", "DANA", "QRIS"}

	// SMSNotification for xendit default notification
	SMSNotification = []string{"sms"}

	// XenditStatusPaid for xendit status paid
	XenditStatusPaid = "PAID"

	// XenditStatusExpired for xendit status expired
	XenditStatusExpired = "EXPIRED"

	// XenditPlatformFee for xendit platform fee
	XenditPlatformFee = 3000.0

	// XenditDisbursementFee for xendit disbursement fee
	XenditDisbursementFee = 5000.0

	// XenditVATPercentage for xendit vat percentage
	XenditVATPercentage = .11
)
